package tfeditor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// DataSourceFilter is an interface which reads Terraform configuration and
// rewrite a given data source, and writes Terraform configuration.
type DataSourceFilter interface {
	// DataSourceFilter reads Terraform configuration and rewrite a given data source,
	// and writes Terraform configuration.
	DataSourceFilter(*tfwrite.File, *tfwrite.DataSource) (*tfwrite.File, error)
}

// MultiDataSourceFilter is a DataSourceFilter implementation which applies
// multiple data source filters to a given data source in sequence.
type MultiDataSourceFilter struct {
	filters []DataSourceFilter
}

var _ DataSourceFilter = (*MultiDataSourceFilter)(nil)

// NewMultiDataSourceFilter creates a new instance of MultiDataSourceFilter.
func NewMultiDataSourceFilter(filters []DataSourceFilter) DataSourceFilter {
	return &MultiDataSourceFilter{
		filters: filters,
	}
}

// DataSourceFilter applies multiple filters to a given data source in sequence.
func (f *MultiDataSourceFilter) DataSourceFilter(inFile *tfwrite.File, dataSource *tfwrite.DataSource) (*tfwrite.File, error) {
	current := inFile
	for _, f := range f.filters {
		next, err := f.DataSourceFilter(current, dataSource)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current, nil
}

// DataSourcesByTypeFilter is a Filter implementation for applying a filter to
// multiple dataSources with a given data source type.
type DataSourcesByTypeFilter struct {
	dataSourceType string
	filter         DataSourceFilter
}

var _ editor.Filter = (*DataSourcesByTypeFilter)(nil)

// NewDataSourcesByTypeFilter creates a new instance of DataSourcesByTypeFilter.
func NewDataSourcesByTypeFilter(dataSourceType string, filter DataSourceFilter) editor.Filter {
	return &DataSourcesByTypeFilter{
		dataSourceType: dataSourceType,
		filter:         filter,
	}
}

// Filter applies a filter to multiple dataSources with a given data source type.
func (f *DataSourcesByTypeFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	current := tfwrite.NewFile(inFile)
	dataSources := current.FindDataSourcesByType(f.dataSourceType)
	for _, dataSource := range dataSources {
		next, err := f.filter.DataSourceFilter(current, dataSource)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current.Raw(), nil
}
