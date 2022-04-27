package schema

import (
	"encoding/json"
	"testing"
)

func TestDictionaryImportID(t *testing.T) {
	cases := []struct {
		desc         string
		importIDMap  map[string]ImportIDFunc
		resourceType string
		resource     string
		ok           bool
		want         string
	}{
		{
			desc: "simple",
			importIDMap: map[string]ImportIDFunc{
				"foo_test1": ImportIDFuncByAttribute("foo1"),
				"foo_test2": ImportIDFuncByAttribute("foo2"),
			},
			resourceType: "foo_test2",
			resource: `
{
  "foo1": "FOO1",
  "foo2": "FOO2",
  "bar": 1,
  "baz": null
}
`,
			ok:   true,
			want: "FOO2",
		},
		{
			desc: "resource type not found",
			importIDMap: map[string]ImportIDFunc{
				"foo_test1": ImportIDFuncByAttribute("foo1"),
				"foo_test2": ImportIDFuncByAttribute("foo2"),
			},
			resourceType: "foo_test3",
			resource: `
{
  "foo1": "FOO1",
  "foo2": "FOO2",
  "bar": 1,
  "baz": null
}
`,
			ok:   false,
			want: "",
		},
		{
			desc: "invalid resource",
			importIDMap: map[string]ImportIDFunc{
				"foo_test1": ImportIDFuncByAttribute("foo1"),
				"foo_test2": ImportIDFuncByAttribute("foo2"),
				"foo_test3": ImportIDFuncByAttribute("bar"),
			},
			resourceType: "foo_test3",
			resource: `
{
  "foo1": "FOO1",
  "foo2": "FOO2",
  "bar": 1,
  "baz": null
}
`,
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

			d := NewDictionary()
			d.RegisterImportIDFuncMap(tc.importIDMap)
			got, err := d.ImportID(tc.resourceType, r)

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
