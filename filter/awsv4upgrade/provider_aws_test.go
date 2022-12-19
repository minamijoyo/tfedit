package awsv4upgrade

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfeditor"
)

// buildTestProviderFilter is a helper function which builds an editor filter for testing.
func buildTestProviderFilter(f tfeditor.ProviderFilterFunc) editor.Filter {
	return tfeditor.NewFileFilter(
		&ProviderAWSFilter{
			filters: []tfeditor.BlockFilter{
				tfeditor.ProviderFilterFunc(f),
			},
		},
	)
}

func TestProviderAWSFilter(t *testing.T) {
	cases := []struct {
		name string
		src  string
		ok   bool
		want string
	}{
		{
			name: "simple",
			src: `
provider "aws" {
  region = "ap-northeast-1"

  access_key                  = "dummy"
  secret_key                  = "dummy"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_region_validation      = true
  skip_requesting_account_id  = true

  # mock endpoints with localstack
  endpoints {
    s3 = "http://localstack:4566"
  }

  s3_force_path_style = true
}
`,
			ok: true,
			want: `
provider "aws" {
  region = "ap-northeast-1"

  access_key                  = "dummy"
  secret_key                  = "dummy"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_region_validation      = true
  skip_requesting_account_id  = true

  # mock endpoints with localstack
  endpoints {
    s3 = "http://localstack:4566"
  }

  s3_use_path_style = true
}
`,
		},
		{
			name: "multiple providers",
			src: `
provider "aws" {
  alias = "foo"

  s3_force_path_style = true
}

provider "aws" {
  alias = "bar"

  s3_force_path_style = true
}
`,
			ok: true,
			want: `
provider "aws" {
  alias = "foo"

  s3_use_path_style = true
}

provider "aws" {
  alias = "bar"

  s3_use_path_style = true
}
`,
		},
		{
			name: "argument not found",
			src: `
provider "aws" {
  alias = "foo"
}
`,
			ok: true,
			want: `
provider "aws" {
  alias = "foo"
}
`,
		},
		{
			name: "provider type not found",
			src: `
provider "google" {
  alias = "foo"
}
`,
			ok: true,
			want: `
provider "google" {
  alias = "foo"
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filter := tfeditor.NewFileFilter(NewProviderAWSFilter())
			o := editor.NewEditOperator(filter)
			output, err := o.Apply([]byte(tc.src), "test")
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := string(output)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}
