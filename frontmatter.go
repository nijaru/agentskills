package agentskills

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ToolList parses the spec's space-delimited allowed-tools format.
// YAML sequences are also accepted for convenience; the canonical
// representation is always the space-delimited string form.
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
	case yaml.SequenceNode:
		tools := make([]string, 0, len(value.Content))
		for _, node := range value.Content {
			tool := strings.TrimSpace(node.Value)
			if tool == "" {
				continue
			}
			tools = append(tools, tool)
		}
		*t = tools
		return nil
	case 0:
		*t = nil
		return nil
	default:
		return fmt.Errorf("allowed-tools must be a string or sequence")
	}
}

func internalMetadata(fields map[string]any) bool {
	if len(fields) == 0 {
		return false
	}
	return truthy(fields["internal"])
}

// truthy handles both Go-native booleans and Python-stringified "true".
func truthy(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(strings.TrimSpace(v), "true")
	default:
		return false
	}
}
