package migration

import (
	"encoding/json"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
)

// Plan is a type which wraps Plan of terraform-json and exposes some
// operations which we need.
type Plan struct {
	raw tfjson.Plan
}

// NewPlan parses a plan file in JSON format and creates a new instance of
// Plan.
func NewPlan(planJSON []byte) (*Plan, error) {
	var raw tfjson.Plan

	if err := json.Unmarshal(planJSON, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse plan file: %s", err)
	}

	plan := &Plan{
		raw: raw,
	}

	return plan, nil
}

// ResourceChanges returns a list of changes in plan.
func (p *Plan) ResourceChanges() []*tfjson.ResourceChange {
	return p.raw.ResourceChanges
}
