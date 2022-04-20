package aws

import "github.com/minamijoyo/tfedit/migration/schema"

func init() {
	registerSchema()
}

// registerSchema defines calculation functions of import ID for each resource type.
// It's expected to be called on initialize by blank import.
func registerSchema() {
	importIDMap := map[string]schema.ImportIDFunc{
		"aws_s3_bucket_accelerate_configuration":             schema.ImportIDFuncByAttribute("bucket"),
		"aws_s3_bucket_acl":                                  schema.ImportIDFuncByMultiAttributes([]string{"bucket", "acl"}, ","),
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
	}

	for k, v := range importIDMap {
		schema.RegisterImportIDFunc(k, v)
	}
}
