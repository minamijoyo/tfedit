package aws

import "github.com/minamijoyo/tfedit/migration/schema"

func init() {
	RegisterSchema()
}

func RegisterSchema() {
	importIDMap := map[string]schema.ImportIDFunc{
		"aws_s3_bucket_acl": func(r schema.Resource) string {
			return r["bucket"].(string) + "," + r["acl"].(string)
		},
	}

	for k, v := range importIDMap {
		schema.RegisterImportIDFunc(k, v)
	}
}
