package webapp

import (
	"net/http"

	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/auth/handler/webapp/viewmodels"
	"github.com/authgear/authgear-server/pkg/auth/webapp"
	"github.com/authgear/authgear-server/pkg/lib/authn/stdattrs"
	"github.com/authgear/authgear-server/pkg/lib/authn/user"
	"github.com/authgear/authgear-server/pkg/lib/session"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/template"
)

var TemplateWebSettingsProfileEditHTML = template.RegisterHTML(
	"web/settings_profile_edit.html",
	components...,
)

func ConfigureSettingsProfileEditRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("OPTION", "GET", "POST").
		WithPathPattern("/settings/profile/:variant/edit")
}

type SettingsProfileEditUserService interface {
	GetRaw(id string) (*user.User, error)
	UpdateStandardAttributes(userID string, stdAttrs map[string]interface{}) error
}

type SettingsProfileEditHandler struct {
	ControllerFactory        ControllerFactory
	BaseViewModel            *viewmodels.BaseViewModeler
	SettingsProfileViewModel *viewmodels.SettingsProfileViewModeler
	Renderer                 Renderer
	Users                    SettingsProfileEditUserService
	ErrorCookie              *webapp.ErrorCookie
}

func (h *SettingsProfileEditHandler) GetData(r *http.Request, rw http.ResponseWriter) (map[string]interface{}, error) {
	userID := session.GetUserID(r.Context())

	data := map[string]interface{}{}

	variant := httproute.GetParam(r, "variant")
	data["Variant"] = variant

	baseViewModel := h.BaseViewModel.ViewModel(r, rw)
	viewmodels.Embed(data, baseViewModel)

	viewModelPtr, err := h.SettingsProfileViewModel.ViewModel(*userID)
	if err != nil {
		return nil, err
	}
	viewmodels.Embed(data, *viewModelPtr)

	return data, nil
}

func (h *SettingsProfileEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctrl, err := h.ControllerFactory.New(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer ctrl.Serve()

	ctrl.Get(func() error {
		data, err := h.GetData(r, w)
		if err != nil {
			return err
		}

		h.Renderer.RenderHTML(w, r, TemplateWebSettingsProfileEditHTML, data)
		return nil
	})

	ctrl.PostAction("save", func() error {
		userID := *session.GetUserID(r.Context())
		m := JSONPointerFormToMap(r.Form)

		writeErr := func(inputErr error) error {
			result := webapp.Result{
				RedirectURI: r.URL.Path,
			}
			cookie, err := h.ErrorCookie.SetError(r, apierrors.AsAPIError(inputErr))
			if err != nil {
				return err
			}
			result.Cookies = append(result.Cookies, cookie)
			result.WriteResponse(w, r)
			return nil
		}

		u, err := h.Users.GetRaw(userID)
		if err != nil {
			return writeErr(err)
		}

		attrs, err := stdattrs.T(u.StandardAttributes).MergedWithJSONPointer(m)
		if err != nil {
			return writeErr(err)
		}

		err = h.Users.UpdateStandardAttributes(userID, attrs)
		if err != nil {
			return writeErr(err)
		}

		result := webapp.Result{RedirectURI: "/settings/profile"}
		result.WriteResponse(w, r)
		return nil
	})
}