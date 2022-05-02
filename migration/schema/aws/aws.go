package aws

import "github.com/minamijoyo/tfedit/migration/schema"

// RegisterSchema defines calculation functions of import ID for each resource type.
func RegisterSchema(d *schema.Dictionary) {
	registerS3Schema(d)
}
