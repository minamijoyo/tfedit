package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// DataSource represents a data block.
// It implements the Block interface.
type DataSource struct {
	*block
}

var _ Block = (*DataSource)(nil)

// NewDataSource creates a new instance of DataSource.
func NewDataSource(block *hclwrite.Block) *DataSource {
	b := newBlock(block)
	return &DataSource{block: b}
}

// NewEmptyDataSource creates a new DataSource with an empty body.
func NewEmptyDataSource(dataSourceType string, dataSourceName string) *DataSource {
	block := hclwrite.NewBlock("data", []string{dataSourceType, dataSourceName})
	return NewDataSource(block)
}

// SchemaType returns a type of data source.
// It returns the first label of block.
// Note that it's not the same as the *hclwrite.Block.Type().
func (r *DataSource) SchemaType() string {
	labels := r.block.raw.Labels()
	return labels[0]
}

// Name returns a name of data source.
// It returns the second label of block.
func (r *DataSource) Name() string {
	labels := r.block.raw.Labels()
	return labels[1]
}

// Count returns a meta argument of count.
// It returns nil if not found.
func (r *DataSource) Count() *Attribute {
	return r.GetAttribute("count")
}

// ForEach returns a meta argument of for_each.
// It returns nil if not found.
func (r *DataSource) ForEach() *Attribute {
	return r.GetAttribute("for_each")
}

// ReferableName returns a name of data source instance which can be referenced
// as a part of address.
// It contains an index reference if count or for_each is set.
// If neither count nor for_each is set, it just returns the name.
func (r *DataSource) ReferableName() string {
	name := r.Name()

	if count := r.Count(); count != nil {
		return name + "[count.index]"
	}

	if forEach := r.ForEach(); forEach != nil {
		return name + "[each.key]"
	}

	return name
}
