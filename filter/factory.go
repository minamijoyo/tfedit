package filter

import (
	"fmt"

	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/filter/awsv4upgrade"
)

// NewFilterByType is a factory method for Filter by type.
func NewFilterByType(filterType string) (editor.Filter, error) {
	switch filterType {
	case "awsv4upgrade":
		return awsv4upgrade.NewAllFilter(), nil
	default:
		return nil, fmt.Errorf("unknown filter type: %s", filterType)
	}
}
