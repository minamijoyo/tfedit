package schema

import (
	"fmt"
)

// Dictionary is a map which defines ImportIDFunc for each resource type.
type Dictionary struct {
	importIDMap map[string]ImportIDFunc
}

// NewDictionary returns a new instance of Dictionary.
func NewDictionary() *Dictionary {
	return &Dictionary{
		importIDMap: make(map[string]ImportIDFunc),
	}
}

// RegisterImportIDFunc registers an ImportIDFunc for a given resource type.
func (d *Dictionary) RegisterImportIDFunc(resourceType string, f ImportIDFunc) {
	d.importIDMap[resourceType] = f
}

// RegisterImportIDFuncMap is a helper method to register a map of ImportIDFunc.
func (d *Dictionary) RegisterImportIDFuncMap(importIDFuncMap map[string]ImportIDFunc) {
	for k, v := range importIDFuncMap {
		d.RegisterImportIDFunc(k, v)
	}
}

// ImportID calculates an import ID from a given resource.
func (d *Dictionary) ImportID(resourceType string, r Resource) (string, error) {
	f, ok := d.importIDMap[resourceType]
	if !ok {
		return "", fmt.Errorf("unknown resource type for import: %s", resourceType)
	}
	return f(r)
}
