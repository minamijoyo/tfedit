package tfeditor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// BlockFilter is an interface which reads Terraform configuration and
// rewrite a given block, and writes Terraform configuration.
type BlockFilter interface {
	// BlockFilter reads Terraform configuration and rewrite a given block,
	// and writes Terraform configuration.
	BlockFilter(*tfwrite.File, tfwrite.Block) (*tfwrite.File, error)
}

type BlockFilterFunc func(*tfwrite.File, tfwrite.Block) (*tfwrite.File, error)

func (f BlockFilterFunc) BlockFilter(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
	return f(inFile, block)
}

// MultiBlockFilter is a BlockFilter implementation which applies
// multiple block filters to a given block in sequence.
type MultiBlockFilter struct {
	filters []BlockFilter
}

var _ BlockFilter = (*MultiBlockFilter)(nil)

// NewMultiBlockFilter creates a new instance of MultiBlockFilter.
func NewMultiBlockFilter(filters []BlockFilter) BlockFilter {
	return &MultiBlockFilter{
		filters: filters,
	}
}

// BlockFilter applies multiple filters to a given block in sequence.
func (f *MultiBlockFilter) BlockFilter(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
	current := inFile
	for _, f := range f.filters {
		next, err := f.BlockFilter(current, block)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current, nil
}

// BlocksByTypeFilter is a Filter implementation for applying a filter to
// multiple blocks with a given block type.
type BlocksByTypeFilter struct {
	blockType  string
	schemaType string
	filter     BlockFilter
}

var _ editor.Filter = (*BlocksByTypeFilter)(nil)

// NewBlocksByTypeFilter creates a new instance of BlocksByTypeFilter.
func NewBlocksByTypeFilter(blockType string, schemaType string, filter BlockFilter) editor.Filter {
	return &BlocksByTypeFilter{
		blockType:  blockType,
		schemaType: schemaType,
		filter:     filter,
	}
}

// Filter applies a filter to multiple blocks with a given block type.
func (f *BlocksByTypeFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	current := tfwrite.NewFile(inFile)
	blocks := current.FindBlocksByType(f.blockType, f.schemaType)
	for _, block := range blocks {
		next, err := f.filter.BlockFilter(current, block)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current.Raw(), nil
}
