package tfwrite

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFileFindResourcesByType(t *testing.T) {
	cases := []struct {
		desc         string
		src          string
		resourceType string
		want         string
		ok           bool
	}{
		{
			desc: "simple",
			src: `
resource "foo_test" "example1" {}
resource "foo_test" "example2" {}
resource "foo_bar" "example1" {}
data "foo_test" "example1" {}
`,
			resourceType: "foo_test",
			want: `
resource "foo_test" "example1" {}

resource "foo_test" "example2" {}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			resources := f.FindResourcesByType(tc.resourceType)

			newFile := NewEmptyFile()
			for _, r := range resources {
				newFile.AppendBlock(r)
			}

			got := printTestFile(t, newFile)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestFileFindProvidersByType(t *testing.T) {
	cases := []struct {
		desc         string
		src          string
		providerType string
		want         string
		ok           bool
	}{
		{
			desc: "simple",
			src: `
provider "aws" {
  region = "ap-northeast-1"
  alias  = "ap_northeast_1"
}

provider "aws" {
  region = "us-east-1"
  alias  = "us_east_1"
}

provider "google" {}

resource "aws_s3_bucket" "example" {}
`,
			providerType: "aws",
			want: `
provider "aws" {
  region = "ap-northeast-1"
  alias  = "ap_northeast_1"
}

provider "aws" {
  region = "us-east-1"
  alias  = "us_east_1"
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			providers := f.FindProvidersByType(tc.providerType)

			newFile := NewEmptyFile()
			for _, p := range providers {
				newFile.AppendBlock(p)
			}

			got := printTestFile(t, newFile)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestFileAppendBlock(t *testing.T) {
	cases := []struct {
		desc         string
		src          string
		resourceType string
		resourceName string
		want         string
		ok           bool
	}{
		{
			desc: "simple",
			src: `
resource "foo_test" "example1" {}
`,
			resourceType: "foo_test",
			resourceName: "example2",
			want: `
resource "foo_test" "example1" {}

resource "foo_test" "example2" {
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			r := NewEmptyResource(tc.resourceType, tc.resourceName)
			f.AppendBlock(r)

			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}
