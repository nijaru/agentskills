package agentskills

import (
	"fmt"
	"strings"
)

// Skill represents a reusable unit of methodology and knowledge.
// It is defined by a SKILL.md file with YAML frontmatter.
type Skill struct {
	Name         string         `yaml:"name"`
	Description  string         `yaml:"description"`
	License      string         `yaml:"license,omitempty"`
	Compatibility string        `yaml:"compatibility,omitempty"`
	Metadata     map[string]any `yaml:"metadata,omitempty"`
	AllowedTools ToolList       `yaml:"allowed-tools,omitempty"`
	Scripts      []string       `yaml:"scripts,omitempty"`
	Instructions string         `yaml:"-"`
	Internal     bool           `yaml:"-"`
	Extra        map[string]any `yaml:",inline"`
}

// Summary returns a one-line summary of the skill for progressive disclosure.
func (s *Skill) Summary() string {
	return s.Name + ": " + s.Description
}

// Validate checks the required spec fields and known limits for a skill.
func (s *Skill) Validate() error {
	if err := ValidateName(s.Name); err != nil {
		return err
	}
	if strings.TrimSpace(s.Description) == "" {
		return fmt.Errorf("skill description is required")
	}
	if len(s.Description) > 1024 {
		return fmt.Errorf("skill description must be 1-1024 characters, got %d", len(s.Description))
	}
	if len(s.Compatibility) > 500 {
		return fmt.Errorf("skill compatibility must be 1-500 characters, got %d", len(s.Compatibility))
	}
	return nil
}
