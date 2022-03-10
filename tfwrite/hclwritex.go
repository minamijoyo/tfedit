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
	ret := []hclwrite.Tokens{}
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
			ret = append(ret, tokens[begin:end])
			foundCBrack = true
			break
		}

		// Find a `,` token
		if tokens[end].Type == hclsyntax.TokenComma {
			ret = append(ret, tokens[begin:end])
			begin = end + 1 // Move the begin cursor after the `,`
		}
	}

	if !foundCBrack {
		return nil // `]` not found
	}

	return ret
}
