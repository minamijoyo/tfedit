package tfeditor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfwrite"
	"github.com/zclconf/go-cty/cty"
)

func TestBlockFilter(t *testing.T) {
	cases := []struct {
		name string
		src  string
		f    BlockFilter
		ok   bool
		want string
	}{
		{
			name: "block",
			src: `
block "foo_bar" "example" {
  baz = "test1"
}
`,
			f: BlockFilterFunc(func(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
				block.SetAttributeValue("baz", cty.StringVal("test2"))
				return inFile, nil
			}),
			ok: true,
			want: `
block "foo_bar" "example" {
  baz = "test2"
}
`,
		},
		{
			name: "resource",
			src: `
resource "foo_bar" "example" {
  baz = "test1"
}
`,
			f: ResourceFilterFunc(func(inFile *tfwrite.File, block *tfwrite.Resource) (*tfwrite.File, error) {
				block.SetAttributeValue("baz", cty.StringVal("test2"))
				return inFile, nil
			}),
			ok: true,
			want: `
resource "foo_bar" "example" {
  baz = "test2"
}
`,
		},
		{
			name: "data",
			src: `
data "foo_bar" "example" {
  baz = "test1"
}
`,
			f: DataSourceFilterFunc(func(inFile *tfwrite.File, block *tfwrite.DataSource) (*tfwrite.File, error) {
				block.SetAttributeValue("baz", cty.StringVal("test2"))
				return inFile, nil
			}),
			ok: true,
			want: `
data "foo_bar" "example" {
  baz = "test2"
}
`,
		},
		{
			name: "provider",
			src: `
provider "foo" {
  baz = "test1"
}
`,
			f: ProviderFilterFunc(func(inFile *tfwrite.File, block *tfwrite.Provider) (*tfwrite.File, error) {
				block.SetAttributeValue("baz", cty.StringVal("test2"))
				return inFile, nil
			}),
			ok: true,
			want: `
provider "foo" {
  baz = "test2"
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filter := NewAllBlocksFilter(tc.f)
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

func TestMultiBlockFilter(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		filters []BlockFilter
		ok      bool
		want    string
	}{
		{
			name: "block",
			src: `
block "foo_bar" "example" {
  baz1 = "test11"
  baz2 = "test12"
}
`,
			filters: []BlockFilter{
				BlockFilterFunc(func(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
					block.SetAttributeValue("baz1", cty.StringVal("test21"))
					return inFile, nil
				}),
				BlockFilterFunc(func(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
					block.SetAttributeValue("baz2", cty.StringVal("test22"))
					return inFile, nil
				}),
			},
			ok: true,
			want: `
block "foo_bar" "example" {
  baz1 = "test21"
  baz2 = "test22"
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mf := NewMultiBlockFilter(tc.filters)
			filter := NewAllBlocksFilter(mf)
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

func TestAllBlocksFilter(t *testing.T) {
	cases := []struct {
		name   string
		src    string
		filter BlockFilter
		ok     bool
		want   string
	}{
		{
			name: "simple",
			src: `
block "foo_bar1" "example11" {
  baz = "test1"
}

block "foo_bar2" "example21" {
  baz = "test1"
}

block "foo_bar1" "example13" {
  baz = "test1"
}
`,
			filter: BlockFilterFunc(func(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
				block.SetAttributeValue("baz", cty.StringVal("test2"))
				return inFile, nil
			}),
			ok: true,
			want: `
block "foo_bar1" "example11" {
  baz = "test2"
}

block "foo_bar2" "example21" {
  baz = "test2"
}

block "foo_bar1" "example13" {
  baz = "test2"
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filter := NewAllBlocksFilter(tc.filter)
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
