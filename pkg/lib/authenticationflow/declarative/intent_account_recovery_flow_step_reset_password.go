package declarative

import (
	"context"
	"fmt"

	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"

	authflow "github.com/authgear/authgear-server/pkg/lib/authenticationflow"
	"github.com/authgear/authgear-server/pkg/lib/config"
)

func init() {
	authflow.RegisterIntent(&IntentAccountRecoveryFlowStepResetPassword{})
}

type IntentAccountRecoveryFlowStepResetPasswordData struct {
	PasswordPolicy *PasswordPolicy `json:"password_policy,omitempty"`
}

var _ authflow.Data = IntentAccountRecoveryFlowStepResetPasswordData{}

func (IntentAccountRecoveryFlowStepResetPasswordData) Data() {}

type IntentAccountRecoveryFlowStepResetPassword struct {
	StepName    string        `json:"step_name,omitempty"`
	JSONPointer jsonpointer.T `json:"json_pointer,omitempty"`
}

var _ authflow.Intent = &IntentAccountRecoveryFlowStepResetPassword{}
var _ authflow.DataOutputer = &IntentAccountRecoveryFlowStepResetPassword{}

func (*IntentAccountRecoveryFlowStepResetPassword) Kind() string {
	return "IntentAccountRecoveryFlowStepResetPassword"
}

func (i *IntentAccountRecoveryFlowStepResetPassword) CanReactTo(ctx context.Context, deps *authflow.Dependencies, flows authflow.Flows) (authflow.InputSchema, error) {
	if len(flows.Nearest.Nodes) == 0 {
		return &InputSchemaTakeNewPassword{
			JSONPointer: i.JSONPointer,
		}, nil
	}
	return nil, authflow.ErrEOF
}

func (i *IntentAccountRecoveryFlowStepResetPassword) ReactTo(ctx context.Context, deps *authflow.Dependencies, flows authflow.Flows, input authflow.Input) (*authflow.Node, error) {
	milestone, ok := authflow.FindMilestone[MilestoneAccountRecoveryCode](flows.Root)
	if !ok {
		return nil, InvalidFlowConfig.New("IntentAccountRecoveryFlowStepResetPassword depends on MilestoneAccountRecoveryCode")
	}
	code := milestone.MilestoneAccountRecoveryCode()

	var inputTakeNewPassword inputTakeNewPassword
	if authflow.AsInput(input, &inputTakeNewPassword) {
		newPassword := inputTakeNewPassword.GetNewPassword()
		return authflow.NewNodeSimple(&NodeDoResetPassword{
			Code:        code,
			NewPassword: newPassword,
		}), nil
	}

	return nil, authflow.ErrIncompatibleInput
}

func (*IntentAccountRecoveryFlowStepResetPassword) step(o config.AuthenticationFlowObject) *config.AuthenticationFlowLoginFlowStep {
	step, ok := o.(*config.AuthenticationFlowLoginFlowStep)
	if !ok {
		panic(fmt.Errorf("flow object is %T", o))
	}

	return step
}

func (i *IntentAccountRecoveryFlowStepResetPassword) OutputData(ctx context.Context, deps *authflow.Dependencies, flows authflow.Flows) (authflow.Data, error) {
	return IntentAccountRecoveryFlowStepResetPasswordData{
		PasswordPolicy: NewPasswordPolicy(
			deps.FeatureConfig.Authenticator,
			deps.Config.Authenticator.Password.Policy,
		),
	}, nil
}
