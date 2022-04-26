package aws

import (
	"testing"

	"github.com/minamijoyo/tfedit/migration/schema"
)

func TestImportIDFuncAWSS3BucketACL(t *testing.T) {
	cases := []struct {
		desc string
		r    schema.Resource
		want string
	}{
		{
			desc: "acl",
			r: map[string]interface{}{
				"acl":                   "private",
				"bucket":                "tfedit-test",
				"expected_bucket_owner": nil,
			},
			want: "tfedit-test,private",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := importIDFuncAWSS3BucketACL(tc.r)

			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}
