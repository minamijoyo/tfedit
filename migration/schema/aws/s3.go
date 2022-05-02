package aws

import (
	"fmt"

	"github.com/minamijoyo/tfedit/migration/schema"
)

func registerS3Schema(d *schema.Dictionary) {
	d.RegisterImportIDFuncMap(map[string]schema.ImportIDFunc{
		"aws_s3_bucket_accelerate_configuration":             schema.ImportIDFuncByAttribute("bucket"),
		"aws_s3_bucket_acl":                                  importIDFuncAWSS3BucketACL,
		"aws_s3_bucket_cors_configuration":                   schema.ImportIDFuncByAttribute("bucket"),
		"aws_s3_bucket_lifecycle_configuration":              schema.ImportIDFuncByAttribute("bucket"),
		"aws_s3_bucket_logging":                              schema.ImportIDFuncByAttribute("bucket"),
		"aws_s3_bucket_object_lock_configuration":            schema.ImportIDFuncByAttribute("bucket"),
		"aws_s3_bucket_policy":                               schema.ImportIDFuncByAttribute("bucket"),
		"aws_s3_bucket_replication_configuration":            schema.ImportIDFuncByAttribute("bucket"),
		"aws_s3_bucket_request_payment_configuration":        schema.ImportIDFuncByAttribute("bucket"),
		"aws_s3_bucket_server_side_encryption_configuration": schema.ImportIDFuncByAttribute("bucket"),
		"aws_s3_bucket_versioning":                           schema.ImportIDFuncByAttribute("bucket"),
		"aws_s3_bucket_website_configuration":                schema.ImportIDFuncByAttribute("bucket"),
	})
}

// importIDFuncAWSS3BucketACL is an implementation of importIDFunc for aws_s3_bucket_acl.
// https://registry.terraform.io/providers/hashicorp%20%20/aws/latest/docs/resources/s3_bucket_acl#import
func importIDFuncAWSS3BucketACL(r schema.Resource) (string, error) {
	// The acl argument conflicts with access_control_policy
	switch {
	case r["acl"] != nil && r["access_control_policy"] == nil: // acl
		return schema.ImportIDFuncByMultiAttributes([]string{"bucket", "acl"}, ",")(r)
	case r["acl"] == nil && r["access_control_policy"] != nil: // grant
		return schema.ImportIDFuncByAttribute("bucket")(r)
	default:
		return "", fmt.Errorf("failed to detect an ID of aws_s3_bucket_acl resource for import: %#v", r)
	}
}
