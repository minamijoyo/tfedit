package schema

import (
	"encoding/json"
	"testing"
)

func TestImportIDFuncByAttribute(t *testing.T) {
	cases := []struct {
		desc     string
		resource string
		key      string
		ok       bool
		want     string
	}{
		{
			desc: "simple",
			resource: `
{
  "foo": "FOO",
  "bar": 1,
  "baz": null
}
`,
			key:  "foo",
			ok:   true,
			want: "FOO",
		},
		{
			desc: "type cast error",
			resource: `
{
  "foo": "FOO",
  "bar": 1,
  "baz": null
}
`,
			key:  "bar",
			ok:   false,
			want: "",
		},
		{
			desc: "found null",
			resource: `
{
  "foo": "FOO",
  "bar": 1,
  "baz": null
}
`,
			key:  "baz",
			ok:   false,
			want: "",
		},
		{
			desc: "not found",
			resource: `
{
  "foo": "FOO",
  "bar": 1,
  "baz": null
}
`,
			key:  "qux",
			ok:   false,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			var r Resource
			if err := json.Unmarshal([]byte(tc.resource), &r); err != nil {
				t.Fatalf("failed to unmarshal json: %s", err)
			}

			got, err := ImportIDFuncByAttribute(tc.key)(r)

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

func TestImportIDFuncByMultiAttributes(t *testing.T) {
	cases := []struct {
		desc     string
		resource string
		keys     []string
		sep      string
		ok       bool
		want     string
	}{
		{
			desc: "simple",
			resource: `
{
  "foo1": "FOO1",
  "foo2": "FOO2",
  "bar": 1,
  "baz": null
}
`,
			keys: []string{"foo1", "foo2"},
			sep:  ",",
			ok:   true,
			want: "FOO1,FOO2",
		},
		{
			desc: "type cast error",
			resource: `
{
  "foo1": "FOO1",
  "foo2": "FOO2",
  "bar": 1,
  "baz": null
}
`,
			keys: []string{"foo1", "bar"},
			sep:  ",",
			ok:   false,
			want: "",
		},
		{
			desc: "found null",
			resource: `
{
  "foo1": "FOO1",
  "foo2": "FOO2",
  "bar": 1,
  "baz": null
}
`,
			keys: []string{"foo1", "baz"},
			sep:  ",",
			ok:   false,
			want: "",
		},
		{
			desc: "not found",
			resource: `
{
  "foo1": "FOO1",
  "foo2": "FOO2",
  "bar": 1,
  "baz": null
}
`,
			keys: []string{"foo1", "qux"},
			sep:  ",",
			ok:   false,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			var r Resource
			if err := json.Unmarshal([]byte(tc.resource), &r); err != nil {
				t.Fatalf("failed to unmarshal json: %s", err)
			}

			got, err := ImportIDFuncByMultiAttributes(tc.keys, tc.sep)(r)

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
