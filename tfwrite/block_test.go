package tfwrite

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func TestBlockType(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "no label",
			src:  `terraform  {}`,
			want: "terraform",
			ok:   true,
		},
		{
			desc: "with a single label",
			src:  `provider "aws" {}`,
			want: "provider",
			ok:   true,
		},
		{
			desc: "with multiple labels",
			src:  `resource "aws_s3_bucket" "example" {}`,
			want: "resource",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			got := b.Type()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestBlockSchemaType(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "no label",
			src:  `terraform  {}`,
			want: "",
			ok:   true,
		},
		{
			desc: "with a single label",
			src:  `provider "aws" {}`,
			want: "aws",
			ok:   true,
		},
		{
			desc: "with multiple labels",
			src:  `resource "aws_s3_bucket" "example" {}`,
			want: "aws_s3_bucket",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			got := b.SchemaType()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestBlockSetType(t *testing.T) {
	cases := []struct {
		desc     string
		src      string
		typeName string
		want     string
		ok       bool
	}{
		{
			desc: "simple",
			src: `
foo {}
`,
			typeName: "bar",
			want: `
bar {}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			b.SetType(tc.typeName)
			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestBlockGetAttribute(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		name string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo {
  bar = "baz"
}
`,
			name: "bar",
			want: `"baz"`,
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			attr := b.GetAttribute(tc.name)
			got, err := attr.ValueAsString()

			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestBlockSetAttributeValue(t *testing.T) {
	cases := []struct {
		desc  string
		src   string
		name  string
		value cty.Value
		want  string
		ok    bool
	}{
		{
			desc: "simple",
			src: `
foo {
  bar = "baz"
}
`,
			name:  "bar",
			value: cty.StringVal("qux"),
			want:  `"qux"`,
			ok:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			b.SetAttributeValue(tc.name, tc.value)
			attr := b.GetAttribute(tc.name)
			got, err := attr.ValueAsString()

			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestBlockSetAttributeRaw(t *testing.T) {
	cases := []struct {
		desc   string
		src    string
		name   string
		tokens hclwrite.Tokens
		want   string
		ok     bool
	}{
		{
			desc: "simple",
			src: `
foo {
  bar = "baz"
}
`,
			name:   "bar",
			tokens: hclwrite.TokensForValue(cty.StringVal("qux")),
			want:   `"qux"`,
			ok:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			b.SetAttributeRaw(tc.name, tc.tokens)
			attr := b.GetAttribute(tc.name)
			got, err := attr.ValueAsString()

			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestBlockAppendAttribute(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		name string
		attr *Attribute
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo {
  bar = "baz"
}
`,
			name: "bar",
			want: `
foo {
  bar = "baz"

  nested {
    bar = "baz"
  }
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			nested := NewEmptyNestedBlock("nested")
			b.AppendNestedBlock(nested)
			attr := b.GetAttribute(tc.name)

			nested.AppendAttribute(attr)

			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestBlockRemoveAttribute(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		name string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo {
  bar = "baz"
  qux = "quux"
}
`,
			name: "bar",
			want: `
foo {
  qux = "quux"
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			b.RemoveAttribute(tc.name)

			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestBlockCopyAttribute(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		name string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo1 {
  bar = "baz"
  qux = "quux"
}

foo2 {
  qux = "quux"
}
`,
			name: "bar",
			want: `
foo1 {
  bar = "baz"
  qux = "quux"
}

foo2 {
  qux = "quux"
  bar = "baz"
}
`,
			ok: true,
		},
		{
			desc: "not found",
			src: `
foo1 {
  qux = "quux"
}

foo2 {
  qux = "quux"
}
`,
			name: "bar",
			want: `
foo1 {
  qux = "quux"
}

foo2 {
  qux = "quux"
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			blocks := findTestBlocks(t, f)
			blocks[1].CopyAttribute(blocks[0], tc.name)

			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestBlockAppendNestedBlock(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		name string
		attr *Attribute
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo {
  bar = "baz"
}
`,
			name: "nested",
			want: `
foo {
  bar = "baz"

  nested {
  }
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			nested := NewEmptyNestedBlock(tc.name)
			b.AppendNestedBlock(nested)

			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestBlockAppendUnwrappedNestedBlockBody(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		name string
		attr *Attribute
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo1 {
  nested {
    bar = "baz"
  }
}

foo2 {
}
`,
			name: "nested",
			want: `
foo1 {
  nested {
    bar = "baz"
  }
}

foo2 {

  bar = "baz"
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			blocks := findTestBlocks(t, f)
			nested := blocks[0].FindNestedBlocksByType("nested")[0]
			blocks[1].AppendUnwrappedNestedBlockBody(nested)

			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestBlockRemoveNestedBlock(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		name string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo {
  bar = "baz"
  nested {
    qux = "quux"
  }
}
`,
			name: "nested",
			want: `
foo {
  bar = "baz"
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			nested := b.FindNestedBlocksByType(tc.name)[0]
			b.RemoveNestedBlock(nested)

			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestBlockFindNestedBlocksByType(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		name string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo1 {
  bar = "baz"
  nested {
    qux = "quux1"
  }

  nested {
    qux = "quux2"
  }
}

foo2 {
}
`,
			name: "nested",
			want: `
foo1 {
  bar = "baz"
  nested {
    qux = "quux1"
  }

  nested {
    qux = "quux2"
  }
}

foo2 {

  nested {
    qux = "quux1"
  }

  nested {
    qux = "quux2"
  }
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			blocks := findTestBlocks(t, f)
			nestedBlocks := blocks[0].FindNestedBlocksByType(tc.name)
			for _, nested := range nestedBlocks {
				blocks[1].AppendNestedBlock(nested)
			}

			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestBlockVerticalFormat(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo {

  bar = "baz"

  nested {

    qux = "quux1"

  }


  nested {


    qux = "quux2"
  }
}
`,
			want: `
foo {
  bar = "baz"

  nested {

    qux = "quux1"

  }

  nested {

    qux = "quux2"
  }
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			b.VerticalFormat()

			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestBlockRenameReference(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		from string
		to   string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo {
  bar = aaa.bbb.ccc
}
`,
			from: "aaa.bbb",
			to:   "xxx.yyy",
			want: `
foo {
  bar = xxx.yyy.ccc
}
`,
			ok: true,
		},
		{
			desc: "multiple attributes",
			src: `
foo {
  bar1 = aaa.bbb.ccc
  bar2 = aaa.bbb.ccc
}
`,
			from: "aaa.bbb",
			to:   "xxx.yyy",
			want: `
foo {
  bar1 = xxx.yyy.ccc
  bar2 = xxx.yyy.ccc
}
`,
			ok: true,
		},
		{
			desc: "multiple nested blocks",
			src: `
foo {
  nested {
    bar1 = aaa.bbb.ccc
  }
  nested {
    bar2 = aaa.bbb.ccc
  }
}
`,
			from: "aaa.bbb",
			to:   "xxx.yyy",
			want: `
foo {
  nested {
    bar1 = xxx.yyy.ccc
  }
  nested {
    bar2 = xxx.yyy.ccc
  }
}
`,
			ok: true,
		},
		{
			desc: "multiple level nested blocks",
			src: `
foo {
  nested {
    bar1 = aaa.bbb.ccc
    nested {
      bar2 = aaa.bbb.ccc
    }
  }
}
`,
			from: "aaa.bbb",
			to:   "xxx.yyy",
			want: `
foo {
  nested {
    bar1 = xxx.yyy.ccc
    nested {
      bar2 = xxx.yyy.ccc
    }
  }
}
`,
			ok: true,
		},
		{
			desc: "not found",
			src: `
foo {
  bar = aaa.bbb.ccc
}
`,
			from: "AAA.BBB",
			to:   "xxx.yyy",
			want: `
foo {
  bar = aaa.bbb.ccc
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			b.RenameReference(tc.from, tc.to)

			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}
