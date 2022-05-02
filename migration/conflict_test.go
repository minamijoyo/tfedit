package migration

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/minamijoyo/tfedit/migration/schema"
)

func TestConflictPlannedActionType(t *testing.T) {
	cases := []struct {
		desc string
		c    *Conflict
		want string
	}{
		{
			desc: "create",
			c: &Conflict{
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
			want: "create",
		},
		{
			desc: "unknown (create before destroy)",
			c: &Conflict{
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
							"delete",
						},
					},
				},
				resolved: false,
			},
			want: "unknown",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.c.PlannedActionType()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestConflictResourceType(t *testing.T) {
	cases := []struct {
		desc string
		c    *Conflict
		want string
	}{
		{
			desc: "simple",
			c: &Conflict{
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
			want: "aws_s3_bucket_acl",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.c.ResourceType()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestConflictAddress(t *testing.T) {
	cases := []struct {
		desc string
		c    *Conflict
		want string
	}{
		{
			desc: "simple",
			c: &Conflict{
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
			want: "aws_s3_bucket_acl.example",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.c.Address()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestConflictResourceAfter(t *testing.T) {
	cases := []struct {
		desc string
		c    *Conflict
		ok   bool
		want schema.Resource
	}{
		{
			desc: "simple",
			c: &Conflict{
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
			ok: true,
			want: schema.Resource(map[string]interface{}{
				"acl":                   "private",
				"bucket":                "tfedit-test",
				"expected_bucket_owner": nil,
			}),
		},
		{
			desc: "type cast error",
			c: &Conflict{
				rc: &tfjson.ResourceChange{
					Change: &tfjson.Change{
						After: nil,
					},
				},
				resolved: false,
			},
			ok:   false,
			want: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := tc.c.ResourceAfter()

			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, got: %#v", got)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%#v\nwant:\n%#v\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}
