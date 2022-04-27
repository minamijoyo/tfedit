package migration

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestStateImportResolver(t *testing.T) {
	cases := []struct {
		desc     string
		s        *Subject
		ok       bool
		resolved bool
		want     []StateAction
	}{
		{
			desc: "simple",
			s: &Subject{
				conflicts: []*Conflict{
					{
						rc: &tfjson.ResourceChange{
							Address:       "aws_s3_bucket_acl.example",
							ModuleAddress: "",
							Mode:          "managed",
							Type:          "aws_s3_bucket_acl",
							Name:          "example",
							Index:         nil,
							ProviderName:  "registry.terraform.io/hashicorp/aws",
							DeposedKey:    "",
							Change: &tfjson.Change{
								Actions: tfjson.Actions{
									"create",
								},
								Before: nil,
								After: map[string]interface{}{
									"acl":                   "private",
									"bucket":                "tfedit-test",
									"expected_bucket_owner": nil,
								},
								AfterUnknown: map[string]interface{}{
									"access_control_policy": true,
									"id":                    true,
								},
								BeforeSensitive: false,
								AfterSensitive: map[string]interface{}{
									"access_control_policy": []interface{}{},
								},
							},
						},
						resolved: false,
					},
				},
			},
			ok:       true,
			resolved: true,
			want: []StateAction{
				&StateImportAction{
					address: "aws_s3_bucket_acl.example",
					id:      "tfedit-test,private",
				},
			},
		},
		{
			desc: "multiple actions",
			s: &Subject{
				conflicts: []*Conflict{
					{
						rc: &tfjson.ResourceChange{
							Address: "aws_s3_bucket_acl.example1",
							Type:    "aws_s3_bucket_acl",
							Change: &tfjson.Change{
								Actions: tfjson.Actions{
									"create",
								},
								Before: nil,
								After: map[string]interface{}{
									"acl":                   "private",
									"bucket":                "tfedit-test1",
									"expected_bucket_owner": nil,
								},
							},
						},
						resolved: false,
					},
					{
						rc: &tfjson.ResourceChange{
							Address: "aws_s3_bucket_acl.example2",
							Type:    "aws_s3_bucket_acl",
							Change: &tfjson.Change{
								Actions: tfjson.Actions{
									"create",
								},
								Before: nil,
								After: map[string]interface{}{
									"acl":                   "private",
									"bucket":                "tfedit-test2",
									"expected_bucket_owner": nil,
								},
							},
						},
						resolved: false,
					},
				},
			},
			ok:       true,
			resolved: true,
			want: []StateAction{
				&StateImportAction{
					address: "aws_s3_bucket_acl.example1",
					id:      "tfedit-test1,private",
				},
				&StateImportAction{
					address: "aws_s3_bucket_acl.example2",
					id:      "tfedit-test2,private",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			d := NewDefaultDictionary()
			r := NewStateImportResolver(d)
			subject, actions, err := r.Resolve(tc.s)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, got: %s", actions)
			}

			if subject.IsResolved() != tc.resolved {
				t.Errorf("unexpected the resolved status of subject. got = %t, but want = %t", subject.IsResolved(), tc.resolved)
			}
			if diff := cmp.Diff(actions, tc.want, cmp.AllowUnexported(StateImportAction{})); diff != "" {
				t.Fatalf("got:\n%#v\nwant:\n%#v\ndiff:\n%s", actions, tc.want, diff)
			}
		})
	}
}
