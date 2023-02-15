package latte

import (
	"context"

	"github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity"
	"github.com/authgear/authgear-server/pkg/lib/workflow"
	"github.com/authgear/authgear-server/pkg/util/validation"
)

func init() {
	workflow.RegisterPrivateIntent(&IntentFindVerifyIdentity{})
}

var IntentFindVerifyIdentitySchema = validation.NewSimpleSchema(`
	{
		"type": "object",
		"additionalProperties": false,
		"properties": {
			"user_id": { "type": "string" }
		},
		"required": ["user_id"]
	}
`)

type IntentFindVerifyIdentity struct {
	UserID string `json:"user_id"`
}

func (*IntentFindVerifyIdentity) Kind() string {
	return "latte.IntentFindVerifyIdentity"
}

func (*IntentFindVerifyIdentity) JSONSchema() *validation.SimpleSchema {
	return IntentFindVerifyIdentitySchema
}

func (*IntentFindVerifyIdentity) CanReactTo(ctx context.Context, deps *workflow.Dependencies, w *workflow.Workflow) ([]workflow.Input, error) {
	if len(w.Nodes) == 0 {
		return []workflow.Input{
			&InputTriggerVerification{},
		}, nil
	}
	return nil, workflow.ErrEOF
}

func (i *IntentFindVerifyIdentity) ReactTo(ctx context.Context, deps *workflow.Dependencies, w *workflow.Workflow, input workflow.Input) (*workflow.Node, error) {
	var trigger inputTriggerVerification

	switch {
	case workflow.AsInput(input, &trigger):
		claimName, claimValue := trigger.ClaimToVerify()
		identities, err := deps.Identities.ListByClaim(claimName, claimValue)
		if err != nil {
			return nil, err
		}

		var iden *identity.Info
		for _, ii := range identities {
			if ii.UserID == i.UserID {
				iden = ii
				break
			}
		}
		if iden == nil {
			// FIXME: define new identity not found error?
			return nil, api.ErrUserNotFound
		}

		return workflow.NewSubWorkflow(&IntentVerifyIdentity{
			IdentityInfo: iden,
			IsFromSignUp: false,
		}), nil
	default:
		return nil, workflow.ErrIncompatibleInput
	}
}

func (*IntentFindVerifyIdentity) GetEffects(ctx context.Context, deps *workflow.Dependencies, w *workflow.Workflow) (effs []workflow.Effect, err error) {
	return nil, nil
}

func (i *IntentFindVerifyIdentity) OutputData(ctx context.Context, deps *workflow.Dependencies, w *workflow.Workflow) (interface{}, error) {
	return nil, nil
}

func (i *IntentFindVerifyIdentity) VerifiedIdentity(w *workflow.Workflow) (*NodeVerifiedIdentity, bool) {
	ws := workflow.FindSubWorkflows[*IntentVerifyIdentity](w)
	if len(ws) == 1 {
		w := ws[0]
		return w.Intent.(*IntentVerifyIdentity).VerifiedIdentity(w)
	}
	return nil, false
}