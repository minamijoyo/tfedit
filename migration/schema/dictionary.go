package schema

import (
	"fmt"
	"strings"
)

// Resource is a type which is equivalent to a type of
// Plan.ResourceChanges[].Change.After in hashicorp/terraform-json,
// but map[string]interface{} is too generic,
// so we give it a friendly alias.
// https://pkg.go.dev/github.com/hashicorp/terraform-json#Change
type Resource map[string]interface{}

// ImportIDFunc is a type of function which calculates an import ID from a given resource.
type ImportIDFunc func(r Resource) (string, error)

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

// ImportIDFuncByAttribute is a helper method to define an ImportIDFunc which
// simply uses a specific single attribute as an import ID.
func ImportIDFuncByAttribute(key string) ImportIDFunc {
	return func(r Resource) (string, error) {
		id, ok := r[key].(string)
		if !ok {
			return "", fmt.Errorf("failed to cast %s to string", key)
		}

		return id, nil
	}
}

// ImportIDFuncByMultiAttributes is a helper method to define an ImportIDFunc which
// joins multiple attributes by a given separater.
func ImportIDFuncByMultiAttributes(keys []string, sep string) ImportIDFunc {
	return func(r Resource) (string, error) {
		elems := []string{}
		for _, key := range keys {
			e, ok := r[key].(string)
			if !ok {
				return "", fmt.Errorf("failed to cast %s to string", key)
			}
			elems = append(elems, e)
		}

		return strings.Join(elems, sep), nil
	}
}
