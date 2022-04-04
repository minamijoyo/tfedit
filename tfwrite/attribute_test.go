package tfwrite

import "testing"

func TestAttributeValueAsString(t *testing.T) {
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

func TestAttributeValueAsTokens(t *testing.T) {
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
			want: ` "baz"`,
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			attr := b.GetAttribute(tc.name)
			tokens := attr.ValueAsTokens()
			got := string(tokens.Bytes())

			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}
