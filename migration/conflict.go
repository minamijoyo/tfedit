package migration

import (
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/minamijoyo/tfedit/migration/schema"
)

// Conflict is a planned resource change.
// It also has a status of whether it has already been resolved.
type Conflict struct {
	// A planned resource change.
	rc *tfjson.ResourceChange
	// A flag indicating that it has already been resolved.
	// The state mv operation reduces two conflicts to a single state migration
	// action, so we need a flag to see if it has already been processed.
	resolved bool
}

// NewConflict returns a new instalce of Conflict.
func NewConflict(rc *tfjson.ResourceChange) *Conflict {
	return &Conflict{
		rc:       rc,
		resolved: false,
	}
}

// MarkAsResolved marks the conflict as resolved.
func (c *Conflict) MarkAsResolved() {
	c.resolved = true
}

// IsResolved return true if the conflict has already been resolved.
func (c *Conflict) IsResolved() bool {
	return c.resolved
}

// PlannedActionType returns a string that represents the type of action.
// Currently some actions that may be included in the plan are not supported.
// It returns "unknown" if not supported.
// The valid values are:
//   - create
//   - unknown
func (c *Conflict) PlannedActionType() string {
	switch {
	case c.rc.Change.Actions.Create():
		return "create"
	default:
		return "unknown"
	}

}

// ResourceType returns a resource type. (e.g. aws_s3_bucket_acl)
func (c *Conflict) ResourceType() string {
	return c.rc.Type
}

// Address returns an absolute address. (e.g. aws_s3_bucket_acl.example)
func (c *Conflict) Address() string {
	return c.rc.Address
}

// ResourceAfter retruns a planned resource after change.
// It doesn't contains attributes known after apply.
func (c *Conflict) ResourceAfter() (schema.Resource, error) {
	after, ok := c.rc.Change.After.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to cast the ResourceChange.Change.After object to Resource: %#v", c.rc.Change.After)
	}

	return schema.Resource(after), nil
}
