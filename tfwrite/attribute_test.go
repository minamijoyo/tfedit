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

func TestAttributeRenameReference(t *testing.T) {
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
			from: "aaa.bbb.ccc",
			to:   "xxx.yyy.zzz",
			want: ` xxx.yyy.zzz`,
			ok:   true,
		},
		{
			desc: "prefix",
			src: `
foo {
  bar = aaa.bbb.ccc
}
`,
			from: "aaa.bbb",
			to:   "xxx.yyy",
			want: ` xxx.yyy.ccc`,
			ok:   true,
		},
		{
			desc: "string template",
			src: `
foo {
  bar = "AAA-${aaa.bbb.ccc}-ZZZ"
}
`,
			from: "aaa.bbb",
			to:   "xxx.yyy",
			want: ` "AAA-${xxx.yyy.ccc}-ZZZ"`,
			ok:   true,
		},
		{
			desc: "multiple references in a string template",
			src: `
foo {
  bar = "AAA-${aaa.bbb.ccc}-ZZZ-${aaa.bbb.ccc}"
}
`,
			from: "aaa.bbb",
			to:   "xxx.yyy",
			want: ` "AAA-${xxx.yyy.ccc}-ZZZ-${xxx.yyy.ccc}"`,
			ok:   true,
		},
		{
			desc: "not found",
			src: `
foo {
  bar = aaa.bbb.ccc
}
`,
			from: "AAA.BBB.CCC",
			to:   "xxx.yyy.zzz",
			want: ` aaa.bbb.ccc`,
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			attr := b.GetAttribute("bar")
			attr.RenameReference(tc.from, tc.to)
			tokens := attr.ValueAsTokens()
			got := string(tokens.Bytes())

			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}
