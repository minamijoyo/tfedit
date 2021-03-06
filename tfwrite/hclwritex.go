package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// SplitTokensAsList parses tokens of a given list and splits it as a list of tokens.
// `["foo", "bar"]` => [`"foo"`, `"bar"`]
// Returns nil if the input cannot be parsed as a list.
func SplitTokensAsList(tokens hclwrite.Tokens) []hclwrite.Tokens {
	// At time of this writing, there is no way for this in hclwrite,
	// so this is a naive implementation.
	tmp := []hclwrite.Tokens{}
	begin := 0
	foundOBrack := false
	for ; begin < len(tokens); begin++ {
		// Find a `[` token
		if tokens[begin].Type == hclsyntax.TokenOBrack {
			foundOBrack = true
			break
		}
	}

	if !foundOBrack {
		return nil // `[` not found
	}

	begin++ // Move the begin cursor after the `[`
	end := begin
	foundCBrack := false
	for ; end < len(tokens); end++ {
		// Find a `]` token
		if tokens[end].Type == hclsyntax.TokenCBrack {
			elm := tokens[begin:end]
			if len(elm) > 0 {
				tmp = append(tmp, elm)
			}
			foundCBrack = true
			break
		}

		// Find a `,` token
		if tokens[end].Type == hclsyntax.TokenComma {
			tmp = append(tmp, tokens[begin:end])
			begin = end + 1 // Move the begin cursor after the `,`
		}
	}

	if !foundCBrack {
		return nil // `]` not found
	}

	ret := []hclwrite.Tokens{}
	for _, t := range tmp {
		r := hclwrite.Tokens{}
		for _, elm := range t {
			// Remove a `\n` (new line) token
			if elm.Type == hclsyntax.TokenNewline {
				continue
			}
			r = append(r, elm)
		}
		if len(r) > 0 {
			ret = append(ret, r)
		}
	}

	return ret
}
