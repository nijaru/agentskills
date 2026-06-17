package agentskills

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// NormalizeName applies NFKC normalization per the Agent Skills spec.
func NormalizeName(name string) string {
	return norm.NFKC.String(name)
}

// ValidateName checks whether name is a valid skills package name.
// Call NormalizeName first if the name comes from user input.
func ValidateName(name string) error {
	if len(name) == 0 || len(name) > 64 {
		return fmt.Errorf("skill name must be 1-64 characters, got %d", len(name))
	}
	if name[0] == '-' || name[len(name)-1] == '-' {
		return fmt.Errorf("skill name %q cannot start or end with a hyphen", name)
	}
	if strings.Contains(name, "--") {
		return fmt.Errorf("skill name %q must not contain consecutive hyphens", name)
	}
	for _, r := range name {
		if r == '-' {
			continue
		}
		if unicode.IsLetter(r) {
			if unicode.ToLower(r) != r {
				return fmt.Errorf("skill name %q must be lowercase", name)
			}
			continue
		}
		if !unicode.IsDigit(r) {
			return fmt.Errorf("skill name %q contains invalid character: %c", name, r)
		}
	}
	return nil
}
