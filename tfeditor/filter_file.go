package tfeditor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// FileFilter is a Filter implementation for applying a filter to all blocks in
// a given file.
type FileFilter struct {
	filter BlockFilter
}

var _ editor.Filter = (*FileFilter)(nil)

// NewFileFilter creates a new instance of NewFileFilter.
func NewFileFilter(filter BlockFilter) editor.Filter {
	return &FileFilter{
		filter: filter,
	}
}

// Filter applies a filter to all blocks in a given file.
func (f *FileFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	current := tfwrite.NewFile(inFile)
	blocks := current.Blocks()
	for _, block := range blocks {
		next, err := f.filter.BlockFilter(current, block)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current.Raw(), nil
}
