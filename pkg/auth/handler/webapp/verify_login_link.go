package webapp

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/authgear/authgear-server/pkg/auth/handler/webapp/viewmodels"
	"github.com/authgear/authgear-server/pkg/auth/webapp"
	"github.com/authgear/authgear-server/pkg/lib/authn/otp"
	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/workflow"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/template"
	"github.com/authgear/authgear-server/pkg/util/validation"
)

var TemplateWebVerifyLoginLinkOTPHTML = template.RegisterHTML(
	"web/verify_login_link.html",
	components...,
)

var VerifyLoginLinkOTPSchema = validation.NewSimpleSchema(`
	{
		"type": "object",
		"properties": {
			"x_oob_otp_code": { "type": "string" }
		},
		"required": ["x_oob_otp_code"]
	}
`)

func ConfigureVerifyLoginLinkOTPRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("OPTIONS", "POST", "GET").
		WithPathPattern("/flows/verify_login_link")
}

type VerifyLoginLinkOTPViewModel struct {
	Code       string
	StateQuery LoginLinkOTPPageQueryState
}

func NewVerifyLoginLinkOTPViewModel(r *http.Request) VerifyLoginLinkOTPViewModel {
	code := r.URL.Query().Get("code")

	return VerifyLoginLinkOTPViewModel{
		Code:       code,
		StateQuery: GetLoginLinkStateFromQuery(r),
	}
}

type WorkflowWebsocketEventStore interface {
	Publish(workflowID string, e workflow.Event) error
}

type VerifyLoginLinkOTPHandler struct {
	LoginLinkOTPCodeService     otp.Service
	GlobalSessionServiceFactory *GlobalSessionServiceFactory
	ControllerFactory           ControllerFactory
	BaseViewModel               *viewmodels.BaseViewModeler
	AuthenticationViewModel     *viewmodels.AuthenticationViewModeler
	Renderer                    Renderer
	WorkflowEvents              WorkflowWebsocketEventStore
}

func (h *VerifyLoginLinkOTPHandler) GetData(r *http.Request, rw http.ResponseWriter) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	baseViewModel := h.BaseViewModel.ViewModel(r, rw)
	viewmodels.Embed(data, NewVerifyLoginLinkOTPViewModel(r))
	viewmodels.Embed(data, baseViewModel)
	return data, nil
}

func (h *VerifyLoginLinkOTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctrl, err := h.ControllerFactory.New(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer ctrl.Serve()

	finishWithState := func(state LoginLinkOTPPageQueryState) {
		url := url.URL{}
		url.Path = r.URL.Path
		query := r.URL.Query()
		query.Set(LoginLinkOTPPageQueryStateKey, string(state))
		url.RawQuery = query.Encode()

		result := webapp.Result{
			RedirectURI:      url.String(),
			NavigationAction: "replace",
		}
		result.WriteResponse(w, r)
	}

	ctrl.Get(func() error {
		data, err := h.GetData(r, w)
		if err != nil {
			return err
		}

		if GetLoginLinkStateFromQuery(r) == LoginLinkOTPPageQueryStateInitial {
			code := r.URL.Query().Get("code")
			_, err := h.LoginLinkOTPCodeService.VerifyLoginLinkCode(code)
			if errors.Is(err, otp.ErrInvalidLoginLink) {
				finishWithState(LoginLinkOTPPageQueryStateInvalidCode)
				return nil
			} else if err != nil {
				return err
			}
		}

		h.Renderer.RenderHTML(w, r, TemplateWebVerifyLoginLinkOTPHTML, data)
		return nil
	})

	ctrl.PostAction("", func() error {
		err := VerifyLoginLinkOTPSchema.Validator().ValidateValue(FormToJSON(r.Form))
		if err != nil {
			return err
		}

		code := r.Form.Get("x_oob_otp_code")

		codeModel, err := h.LoginLinkOTPCodeService.SetUserInputtedLoginLinkCode(code)
		if errors.Is(err, otp.ErrInvalidLoginLink) {
			finishWithState(LoginLinkOTPPageQueryStateInvalidCode)
			return nil
		} else if err != nil {
			return err
		}

		// Update the web session and trigger the refresh event
		if codeModel.WebSessionID != "" {
			webSessionProvider := h.GlobalSessionServiceFactory.NewGlobalSessionService(
				config.AppID(codeModel.AppID),
			)
			webSession, err := webSessionProvider.GetSession(codeModel.WebSessionID)
			if err != nil {
				return err
			}
			err = webSessionProvider.UpdateSession(webSession)
			if err != nil {
				return err
			}
		}

		if codeModel.WorkflowID != "" {
			err = h.WorkflowEvents.Publish(codeModel.WorkflowID, workflow.NewEventRefresh())
			if err != nil {
				return err
			}
		}

		finishWithState(LoginLinkOTPPageQueryStateMatched)
		return nil
	})
}