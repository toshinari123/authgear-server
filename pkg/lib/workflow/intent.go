package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

type Intent interface {
	Kind() string
	Instantiate(data json.RawMessage) error
	GetEffects(ctx context.Context, deps *Dependencies, workflow *Workflow) (effs []Effect, err error)
	DeriveEdges(ctx context.Context, deps *Dependencies, workflow *Workflow) ([]Edge, error)
	OutputData(ctx context.Context, deps *Dependencies, workflow *Workflow) (interface{}, error)
}

type IntentOutput struct {
	Kind string      `json:"kind"`
	Data interface{} `json:"data,omitempty"`
}

type IntentJSON struct {
	Kind string          `json:"kind"`
	Data json.RawMessage `json:"data"`
}

type IntentFactory func() Intent

var intentRegistry = map[string]IntentFactory{}

func RegisterIntent(intent Intent) {
	intentType := reflect.TypeOf(intent).Elem()

	intentKind := intent.Kind()
	factory := IntentFactory(func() Intent {
		return reflect.New(intentType).Interface().(Intent)
	})

	if _, hasKind := intentRegistry[intentKind]; hasKind {
		panic("interaction: duplicated intent kind: " + intentKind)
	}
	intentRegistry[intentKind] = factory
}

func InstantiateIntent(j IntentJSON) (Intent, error) {
	factory, ok := intentRegistry[j.Kind]
	if !ok {
		return nil, fmt.Errorf("workflow: unknown intent kind: %v", j.Kind)
	}
	intent := factory()

	err := intent.Instantiate(j.Data)
	if err != nil {
		return nil, err
	}

	return intent, nil
}
