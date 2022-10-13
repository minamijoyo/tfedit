package tfwrite

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// parseTestFile is a helper for parsing a test file.
func parseTestFile(t *testing.T, src string) *File {
	t.Helper()
	f, diags := hclwrite.ParseConfig([]byte(src), "", hcl.Pos{Line: 1, Column: 1})
	if len(diags) != 0 {
		for _, diag := range diags {
			t.Logf("- %s", diag.Error())
		}
		t.Fatalf("unexpected diagnostics")
	}

	return NewFile(f)
}

// printTestFile is a helper for print a test file.
func printTestFile(t *testing.T, f *File) string {
	t.Helper()
	bytes := f.Raw().BuildTokens(nil).Bytes()
	return string(hclwrite.Format(bytes))
}

// findFirstTestBlock is a test helper for find the first block.
func findFirstTestBlock(t *testing.T, f *File) *block {
	t.Helper()
	blocks := f.Raw().Body().Blocks()
	return newBlock(blocks[0])
}

// findTestBlocks is a test helper for find blocks.
func findTestBlocks(t *testing.T, f *File) []*block {
	t.Helper()
	var blocks []*block
	for _, b := range f.Raw().Body().Blocks() {
		blocks = append(blocks, newBlock(b))
	}
	return blocks
}

// findFirstTestResource is a test helper for find the first resource.
func findFirstTestResource(t *testing.T, f *File) *Resource {
	t.Helper()
	for _, block := range f.raw.Body().Blocks() {
		if block.Type() == "resource" {
			return NewResource(block)
		}
	}
	return nil
}

// findFirstTestDataSource is a test helper for find the first data source.
func findFirstTestDataSource(t *testing.T, f *File) *DataSource {
	t.Helper()
	for _, block := range f.raw.Body().Blocks() {
		if block.Type() == "data" {
			return NewDataSource(block)
		}
	}
	return nil
}

// findFirstTestProvider is a test helper for find the first provider.
func findFirstTestProvider(t *testing.T, f *File) *Provider {
	t.Helper()
	for _, block := range f.raw.Body().Blocks() {
		if block.Type() == "provider" {
			return NewProvider(block)
		}
	}
	return nil
}
