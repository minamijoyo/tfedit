package tfwrite

import (
	"testing"
)

func TestDataSourceType(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
data "foo_test" "example" {}
`,
			want: "data",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			r := findFirstTestDataSource(t, f)

			got := r.Type()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestDataSourceSchemaType(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
data "foo_test" "example" {}
`,
			want: "foo_test",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			r := findFirstTestDataSource(t, f)
			got := r.SchemaType()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestDataSourceName(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
data "foo_test" "example" {}
`,
			want: "example",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			r := findFirstTestDataSource(t, f)
			got := r.Name()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestDataSourceReferableName(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "count",
			src: `
data "foo_test" "example" {
  count = 2
}
`,
			want: "example[count.index]",
			ok:   true,
		},
		{
			desc: "for_each",
			src: `
data "foo_test" "example" {
  for_each = toset(["foo", "bar"])
}
`,
			want: "example[each.key]",
			ok:   true,
		},
		{
			desc: "default",
			src: `
data "foo_test" "example" {
}
`,
			want: "example",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			r := findFirstTestDataSource(t, f)
			got := r.ReferableName()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}
