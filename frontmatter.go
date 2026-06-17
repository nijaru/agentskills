package agentskills

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ToolList parses the spec's space-delimited allowed-tools format.
type ToolList []string

func (t *ToolList) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		if strings.TrimSpace(value.Value) == "" {
			*t = nil
			return nil
		}
		*t = strings.Fields(value.Value)
		return nil
	case 0:
		*t = nil
		return nil
	default:
		return fmt.Errorf("allowed-tools must be a space-delimited string")
	}
}

func internalMetadata(fields map[string]string) bool {
	if len(fields) == 0 {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(fields["internal"]), "true")
}
