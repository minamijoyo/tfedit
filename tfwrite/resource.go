package tfwrite

import "github.com/hashicorp/hcl/v2/hclwrite"

// FindResourcesByType returns all matching blocks from the body that have the
// given resourceType or returns an empty list if there is no matching block.
// This method is useful when you want to ignore the resource name.
func FindResourcesByType(body *hclwrite.Body, resourceType string) []*hclwrite.Block {
	var matched []*hclwrite.Block

	for _, block := range body.Blocks() {
		if block.Type() != "resource" {
			continue
		}

		labels := block.Labels()
		if len(labels) == 2 && labels[0] != resourceType {
			continue
		}

		matched = append(matched, block)
	}

	return matched
}

// GetResourceName is a helper method for getting a resource name of the given block.
func GetResourceName(block *hclwrite.Block) string {
	labels := block.Labels()
	return labels[1]
}

// AppendNewResource is a helper method for appending a new resource block
// to the given body and returns a new block.
func AppendNewResource(body *hclwrite.Body, resourceType string, resourceName string) *hclwrite.Block {
	body.AppendNewline()
	return body.AppendNewBlock("resource", []string{resourceType, resourceName})
}
