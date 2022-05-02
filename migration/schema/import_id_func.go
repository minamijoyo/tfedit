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

// ImportIDFuncByAttribute is a helper method to define an ImportIDFunc which
// simply uses a specific single attribute as an import ID.
func ImportIDFuncByAttribute(key string) ImportIDFunc {
	return func(r Resource) (string, error) {
		id, ok := r[key].(string)
		if !ok {
			return "", fmt.Errorf("failed to cast %s = %#v to string as import ID", key, r[key])
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
				return "", fmt.Errorf("failed to cast %s = %#v to string as an element of import ID", key, r[key])
			}
			elems = append(elems, e)
		}

		return strings.Join(elems, sep), nil
	}
}
