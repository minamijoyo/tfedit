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
		ok       bool
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
			ok:   true,
			want: "tfedit-test,private",
		},
		{
			desc: "grant",
			resource: `
{
  "access_control_policy": [
    {
      "grant": [
        {
          "grantee": [
            {
              "email_address": "",
              "id": "",
              "type": "Group",
              "uri": "http://acs.amazonaws.com/groups/s3/LogDelivery"
            }
          ],
          "permission": "READ_ACP"
        },
        {
          "grantee": [
            {
              "email_address": "",
              "id": "",
              "type": "Group",
              "uri": "http://acs.amazonaws.com/groups/s3/LogDelivery"
            }
          ],
          "permission": "WRITE"
        },
        {
          "grantee": [
            {
              "email_address": "",
              "id": "bcaf1ffd86f41161ca5fb16fd081034f",
              "type": "CanonicalUser",
              "uri": ""
            }
          ],
          "permission": "FULL_CONTROL"
        }
      ],
      "owner": [
        {
          "id": "set_aws_canonical_user_id"
        }
      ]
    }
  ],
  "acl": null,
  "bucket": "tfedit-test",
  "expected_bucket_owner": null
}
`,
			ok:   true,
			want: "tfedit-test",
		},
		{
			desc: "invalid",
			resource: `
{
  "bucket": "tfedit-test",
  "expected_bucket_owner": null
}
`,
			ok:   false,
			want: "",
		},
		{
			desc: "conflict",
			resource: `
{
  "acl": "private",
  "bucket": "tfedit-test",
  "expected_bucket_owner": null,
  "access_control_policy": []
}
`,
			ok:   false,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			var r schema.Resource
			if err := json.Unmarshal([]byte(tc.resource), &r); err != nil {
				t.Fatalf("failed to unmarshal json: %s", err)
			}

			got, err := importIDFuncAWSS3BucketACL(r)

			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, got: %s", got)
			}

			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}
