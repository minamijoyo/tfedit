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

// BlockFilterFunc is a helper method for implementing a BlockFilter interface.
type BlockFilterFunc func(*tfwrite.File, tfwrite.Block) (*tfwrite.File, error)

// BlockFilter reads Terraform configuration and rewrite a given block,
// and writes Terraform configuration.
func (f BlockFilterFunc) BlockFilter(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
	return f(inFile, block)
}

// anyBlockFilterFunc is a generic helper method for implementing a BlockFilter
// interface.
// It would be great if we could treat ResourceFilter and other block filters
// together in an orthogonal way and, at the same time, provide derived types
// such as Resource to library users. To do this, We start using Generics here
// to avoid code duplication, but it is not runtime safe because of the
// interface cast. We reached this design decision after balancing the
// convenience of library users and maintainers. However, since we are not
// confident that the design is correct, we will only expose the derived type
// aliases, and the generic implementation is kept private.
type anyBlockFilterFunc[T tfwrite.Block] func(*tfwrite.File, T) (*tfwrite.File, error)

// ResourceFilterFunc is a helper method for implementing a BlockFilter interface for Resource.
type ResourceFilterFunc = anyBlockFilterFunc[*tfwrite.Resource]

// DataSourceFilterFunc is a helper method for implementing a BlockFilter interface for DataSource.
type DataSourceFilterFunc = anyBlockFilterFunc[*tfwrite.DataSource]

// ProviderFilterFunc is a helper method for implementing a BlockFilter interface for Provider.
type ProviderFilterFunc = anyBlockFilterFunc[*tfwrite.Provider]

// BlockFilter reads Terraform configuration and rewrite a given block,
// and writes Terraform configuration.
func (f anyBlockFilterFunc[T]) BlockFilter(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
	dereived, ok := block.(T)
	if !ok {
		// If the type does not match, it is simply ignored without error.
		// This design is based on the assumption that multiple types of block
		// filters are to be applied to all blocks at once.
		return inFile, nil
	}
	return f(inFile, dereived)
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

// AllBlocksFilter is a Filter implementation for applying a filter to all
// blocks in a given file.
type AllBlocksFilter struct {
	filter BlockFilter
}

var _ editor.Filter = (*AllBlocksFilter)(nil)

// NewAllBlocksFilter creates a new instance of NewAllBlocksFilter.
func NewAllBlocksFilter(filter BlockFilter) editor.Filter {
	return &BlocksByTypeFilter{
		filter: filter,
	}
}

// Filter applies a filter to all blocks in a given file.
func (f *AllBlocksFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	bf := NewBlocksByTypeFilter("", "", f.filter)
	return bf.Filter(inFile)
}
