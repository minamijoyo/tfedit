package awsv4upgrade

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfeditor"
)

func TestProviderAWSS3ForcePathStyleFilter(t *testing.T) {
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
  s3_force_path_style = true
}
`,
			ok: true,
			want: `
provider "aws" {
  s3_use_path_style = true
}
`,
		},
		{
			name: "argument not found",
			src: `
provider "aws" {
  region = "ap-northeast-1"
}
`,
			ok: true,
			want: `
provider "aws" {
  region = "ap-northeast-1"
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filter := &ProviderAWSFilter{filters: []tfeditor.ProviderFilter{&ProviderAWSS3ForcePathStyleFilter{}}}
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
