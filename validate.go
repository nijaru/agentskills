package agentskills

import (
	"fmt"
	"regexp"
	"strings"
)

var skillNameRe = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// ValidateName checks whether name is a valid skills package name.
func ValidateName(name string) error {
	if len(name) == 0 || len(name) > 64 {
		return fmt.Errorf("skill name must be 1-64 characters, got %d", len(name))
	}
	if !skillNameRe.MatchString(name) {
		return fmt.Errorf("skill name %q must match ^[a-z0-9]([a-z0-9-]*[a-z0-9])?$", name)
	}
	if strings.Contains(name, "--") {
		return fmt.Errorf("skill name %q must not contain consecutive hyphens", name)
	}
	return nil
}
