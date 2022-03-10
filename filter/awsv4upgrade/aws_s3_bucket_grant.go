package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
	"github.com/zclconf/go-cty/cty"
)

// AWSS3BucketGrantFilter is a filter implementation for upgrading the
// grant argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#grant-argument
type AWSS3BucketGrantFilter struct{}

var _ tfeditor.ResourceFilter = (*AWSS3BucketGrantFilter)(nil)

// NewAWSS3BucketGrantFilter creates a new instance of AWSS3BucketGrantFilter.
func NewAWSS3BucketGrantFilter() tfeditor.ResourceFilter {
	return &AWSS3BucketGrantFilter{}
}

// ResourceFilter upgrades the grant argument of aws_s3_bucket.
func (f *AWSS3BucketGrantFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "grant"
	newResourceType := "aws_s3_bucket_acl"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setBucketArgument(newResource, resource)

	acpBlock := tfwrite.NewEmptyNestedBlock("access_control_policy")
	newResource.AppendNestedBlock(acpBlock)

	for _, nestedBlock := range nestedBlocks {
		// Split permissions to each grant block
		// A permissions attribute of grant block was a list in v3,
		// but in v4 we need to set each permission to each grant block respectively.
		// grant {
		//   type        = "Group"
		//   permissions = ["READ_ACP", "WRITE"]
		//   uri         = "http://acs.amazonaws.com/groups/s3/LogDelivery"
		// }
		// =>
		// grant {
		//   grantee {
		//     type = "Group"
		//     uri  = "http://acs.amazonaws.com/groups/s3/LogDelivery"
		//   }
		//   permission = "READ_ACP"
		// }
		//
		// grant {
		//   grantee {
		//     type = "Group"
		//     uri  = "http://acs.amazonaws.com/groups/s3/LogDelivery"
		//   }
		//   permission = "WRITE"
		//}
		permissionsAttr := nestedBlock.GetAttribute("permissions")
		if permissionsAttr == nil {
			// The `permissions` attrubute is required, skip if not found.
			continue
		}

		// `["READ_ACP", "WRITE"]` => ["READ_ACP", "WRITE"]
		permissions := tfwrite.SplitTokensAsList(permissionsAttr.ValueAsTokens())
		if permissions == nil {
			// The `permissions` attrubute cannot be parsed as a list.
			// If the `permissions` attribute is passed as a variable or generated by a function,
			// it cannot be split automatically.
			continue
		}

		for _, permission := range permissions {
			grantBlock := tfwrite.NewEmptyNestedBlock("grant")
			acpBlock.AppendNestedBlock(grantBlock)
			granteeBlock := tfwrite.NewEmptyNestedBlock("grantee")
			grantBlock.AppendNestedBlock(granteeBlock)
			nestedBlock.RemoveAttribute("permissions")
			granteeBlock.AppendUnwrappedNestedBlockBody(nestedBlock)
			grantBlock.SetAttributeRaw("permission", permission)
		}

		resource.RemoveNestedBlock(nestedBlock)
	}

	ownerBlock := tfwrite.NewEmptyNestedBlock("owner")
	acpBlock.AppendNestedBlock(ownerBlock)
	// A grant argument of aws_s3_bucket in v3 doesn’t have an owner block,
	// https://registry.terraform.io/providers/hashicorp/aws/3.74.3/docs/resources/s3_bucket#grant
	// but an access_control_policy argument of aws_s3_bucket_acl in v4 has an owner block as required.
	// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_acl#access_control_policy
	// There is no way to set it automatically without the AWS API call.
	ownerBlock.SetAttributeValue("id", cty.StringVal("set_aws_canonical_user_id"))

	return inFile, nil
}
