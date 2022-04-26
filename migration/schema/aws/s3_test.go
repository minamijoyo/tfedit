package aws

import (
	"encoding/json"
	"testing"

	"github.com/minamijoyo/tfedit/migration/schema"
)

func TestImportIDFuncAWSS3BucketACL(t *testing.T) {
	cases := []struct {
		desc     string
		resource string
		want     string
	}{
		{
			desc: "acl",
			resource: `
{
  "acl": "private",
  "bucket": "tfedit-test",
  "expected_bucket_owner": null
}
`,
			want: "tfedit-test,private",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			var r schema.Resource
			if err := json.Unmarshal([]byte(tc.resource), &r); err != nil {
				t.Fatalf("failed to unmarshal json: %s", err)
			}

			got := importIDFuncAWSS3BucketACL(r)

			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}
