package latte

import (
	"context"

	"github.com/authgear/authgear-server/pkg/api/model"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator"
	"github.com/authgear/authgear-server/pkg/lib/workflow"
	"github.com/authgear/authgear-server/pkg/util/uuid"
	"github.com/authgear/authgear-server/pkg/util/validation"
)

func init() {
	workflow.RegisterPrivateIntent(&IntentCreatePassword{})
}

var IntentCreatePasswordSchema = validation.NewSimpleSchema(`{}`)

type IntentCreatePassword struct {
	UserID                 string             `json:"user_id"`
	AuthenticatorKind      authenticator.Kind `json:"authenticator_kind"`
	AuthenticatorIsDefault bool               `json:"authenticator_is_default"`
}

func (*IntentCreatePassword) Kind() string {
	return "latte.IntentCreatePassword"
}

func (*IntentCreatePassword) JSONSchema() *validation.SimpleSchema {
	return IntentCreatePasswordSchema
}

func (*IntentCreatePassword) CanReactTo(ctx context.Context, deps *workflow.Dependencies, w *workflow.Workflow) ([]workflow.Input, error) {
	if len(w.Nodes) == 0 {
		return []workflow.Input{
			&InputCreatePassword{},
		}, nil
	}
	return nil, workflow.ErrEOF
}

func (i *IntentCreatePassword) ReactTo(ctx context.Context, deps *workflow.Dependencies, w *workflow.Workflow, input workflow.Input) (*workflow.Node, error) {
	var inputCreatePassword inputCreatePassword

	switch {
	case workflow.AsInput(input, &inputCreatePassword):
		spec := &authenticator.Spec{
			UserID:    i.UserID,
			IsDefault: i.AuthenticatorIsDefault,
			Kind:      i.AuthenticatorKind,
			Type:      model.AuthenticatorTypePassword,
			Password: &authenticator.PasswordSpec{
				PlainPassword: inputCreatePassword.GetNewPassword(),
			},
		}

		authenticatorID := uuid.New()

		info, err := deps.Authenticators.NewWithAuthenticatorID(authenticatorID, spec)
		if err != nil {
			return nil, err
		}

		return workflow.NewNodeSimple(&NodeDoCreateAuthenticator{
			AuthenticatorInfo: info,
		}), nil
	default:
		return nil, workflow.ErrIncompatibleInput
	}

}

func (*IntentCreatePassword) GetEffects(ctx context.Context, deps *workflow.Dependencies, w *workflow.Workflow) (effs []workflow.Effect, err error) {
	return nil, nil
}

func (*IntentCreatePassword) OutputData(ctx context.Context, deps *workflow.Dependencies, w *workflow.Workflow) (interface{}, error) {
	return nil, nil
}