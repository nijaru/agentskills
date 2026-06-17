package agentskills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoader(t *testing.T) {
	content := `---
name: test-skill
description: A test skill for testing.
allowed-tools: [bash, read_file]
license: Apache-2.0
compatibility: Requires git and jq
metadata:
  author: example-org
  version: "1.0"
---
# Instructions
Use this skill for testing purposes.
`
	tmp := t.TempDir()
	path := filepath.Join(tmp, "SKILL.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if s.Name != "test-skill" {
		t.Errorf("expected name test-skill, got %s", s.Name)
	}
	if s.Description != "A test skill for testing." {
		t.Errorf("expected description, got %s", s.Description)
	}
	if s.License != "Apache-2.0" {
		t.Errorf("expected license Apache-2.0, got %s", s.License)
	}
	if s.Compatibility != "Requires git and jq" {
		t.Errorf("expected compatibility, got %s", s.Compatibility)
	}
	if len(s.AllowedTools) != 2 || s.AllowedTools[0] != "bash" {
		t.Errorf("expected allowed-tools [bash read_file], got %v", s.AllowedTools)
	}
	if s.Metadata["author"] != "example-org" {
		t.Errorf("expected metadata.author example-org, got %v", s.Metadata["author"])
	}
	if s.Metadata["version"] != "1.0" {
		t.Errorf("expected metadata.version 1.0, got %v", s.Metadata["version"])
	}
	if s.Instructions != "# Instructions\nUse this skill for testing purposes." {
		t.Errorf("expected instructions, got %s", s.Instructions)
	}
}

func TestLoaderScalarAllowedToolsAndInternal(t *testing.T) {
	content := `---
name: internal-skill
description: Internal test skill.
allowed-tools: Bash(git:*) Bash(jq:*) Read
metadata:
  internal: true
  owner: example
---
Body.
`
	tmp := t.TempDir()
	path := filepath.Join(tmp, "SKILL.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Internal {
		t.Fatal("expected internal skill to be marked internal")
	}
	if len(s.AllowedTools) != 3 || s.AllowedTools[0] != "Bash(git:*)" {
		t.Fatalf("unexpected allowed-tools: %v", s.AllowedTools)
	}
}

func TestRegisterDeregister(t *testing.T) {
	reg := NewRegistry()
	s := &Skill{Name: "my-skill", Description: "desc"}
	reg.Register(s)
	got, ok := reg.Get("my-skill")
	if !ok || got != s {
		t.Fatal("skill not found after Register")
	}
	reg.Deregister("my-skill")
	if _, ok := reg.Get("my-skill"); ok {
		t.Fatal("skill still present after Deregister")
	}
	reg.Deregister("my-skill")
}

func TestValidateName(t *testing.T) {
	cases := []struct {
		name string
		ok   bool
	}{
		{"hello", true},
		{"my-skill", true},
		{"a1b2", true},
		{"a", true},
		{"", false},
		{"A-skill", false},
		{"-start", false},
		{"end-", false},
		{"no--double", false},
		{"has space", false},
		{"my_skill", false},
		{strings.Repeat("a", 64), true},
		{strings.Repeat("a", 65), false},
		{"技能", true},
		{"мой-навык", true},
		{"навык", true},
		{"НАВЫК", false},
	}
	for _, c := range cases {
		err := ValidateName(c.name)
		if c.ok && err != nil {
			t.Errorf("name=%q: unexpected error: %v", c.name, err)
		}
		if !c.ok && err == nil {
			t.Errorf("name=%q: expected error", c.name)
		}
	}
}

func TestValidateNameNFKCNormalization(t *testing.T) {
	// café with combining accent (5 chars) should normalize to precomposed (4 chars)
	decomposed := "cafe\u0301"
	composed := "caf\u00e9"

	// Both should be valid after normalization
	if err := ValidateName(NormalizeName(decomposed)); err != nil {
		t.Errorf("decomposed form: unexpected error: %v", err)
	}
	if err := ValidateName(NormalizeName(composed)); err != nil {
		t.Errorf("composed form: unexpected error: %v", err)
	}
}

func TestLoaderNormalizesName(t *testing.T) {
	// Decomposed café should be stored as normalized (precomposed) form
	decomposed := "cafe\u0301"
	composed := "caf\u00e9"

	content := fmt.Sprintf("---\nname: %s\ndescription: Test skill.\n---\nBody.\n", decomposed)
	tmp := t.TempDir()
	path := filepath.Join(tmp, "SKILL.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != composed {
		t.Errorf("expected normalized name %q, got %q", composed, s.Name)
	}
}

func TestValidateSkill(t *testing.T) {
	s := &Skill{Name: "pdf-processing", Description: "Extract PDFs"}
	if err := s.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}

	if err := (&Skill{Name: "bad name", Description: "x"}).Validate(); err == nil {
		t.Fatal("expected invalid name to fail")
	}
	if err := (&Skill{Name: "ok", Description: ""}).Validate(); err == nil {
		t.Fatal("expected empty description to fail")
	}
}

func TestValidateDescriptionTooLong(t *testing.T) {
	s := &Skill{
		Name:        "test",
		Description: strings.Repeat("x", 1025),
	}
	if err := s.Validate(); err == nil {
		t.Fatal("expected description too long to fail")
	}
}

func TestValidateCompatibilityTooLong(t *testing.T) {
	s := &Skill{
		Name:         "test",
		Description:  "desc",
		Compatibility: strings.Repeat("x", 501),
	}
	if err := s.Validate(); err == nil {
		t.Fatal("expected compatibility too long to fail")
	}
}

func TestValidateRejectsUnknownFields(t *testing.T) {
	s := &Skill{
		Name:        "test",
		Description: "desc",
		Extra:       map[string]any{"typo-field": "value", "another": "val"},
	}
	err := s.Validate()
	if err == nil {
		t.Fatal("expected unknown fields to be rejected")
	}
	msg := err.Error()
	if !strings.Contains(msg, "typo-field") || !strings.Contains(msg, "another") {
		t.Fatalf("error should list both unknown fields, got: %s", msg)
	}
}

func TestLoaderMissingFrontmatter(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "SKILL.md")
	if err := os.WriteFile(path, []byte("# No frontmatter here"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected missing frontmatter to fail")
	}
}

func TestLoaderUnclosedFrontmatter(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "SKILL.md")
	content := "---\nname: test\ndescription: desc\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected unclosed frontmatter to fail")
	}
}

func TestLoaderInvalidYAML(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "SKILL.md")
	content := "---\nname: [invalid\ndescription: broken\n---\nBody\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected invalid YAML to fail")
	}
}

func TestLoaderNonDictFrontmatter(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "SKILL.md")
	content := "---\n- just\n- a\n- list\n---\nBody\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected YAML list frontmatter to fail")
	}
}

func TestLoaderMissingName(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "SKILL.md")
	content := "---\ndescription: A test skill\n---\nBody\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected missing name to fail")
	}
}

func TestLoaderMissingDescription(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "SKILL.md")
	content := "---\nname: my-skill\n---\nBody\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected missing description to fail")
	}
}

func TestRegistry(t *testing.T) {
	tmp := t.TempDir()
	skillDir := filepath.Join(tmp, "skills", "test")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := `---
name: registry-test
description: Testing registry discovery.
---
Instructions here.
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	reg := NewRegistry(tmp)
	if err := reg.Discover(); err != nil {
		t.Fatal(err)
	}

	skills := reg.List()
	if len(skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(skills))
	}

	s, ok := reg.Get("registry-test")
	if !ok {
		t.Fatal("skill not found in registry")
	}
	if s.Name != "registry-test" {
		t.Errorf("expected name registry-test, got %s", s.Name)
	}
}

func TestRegistryListSorted(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&Skill{Name: "zeta", Description: "last"})
	reg.Register(&Skill{Name: "alpha", Description: "first"})
	reg.Register(&Skill{Name: "mid", Description: "middle"})

	skills := reg.List()
	if len(skills) != 3 {
		t.Fatalf("len(skills) = %d, want 3", len(skills))
	}
	if skills[0].Name != "alpha" || skills[1].Name != "mid" || skills[2].Name != "zeta" {
		t.Fatalf(
			"skill order = [%s %s %s], want [alpha mid zeta]",
			skills[0].Name,
			skills[1].Name,
			skills[2].Name,
		)
	}
}

func TestRegistryDiscoversLowercaseSkillMd(t *testing.T) {
	tmp := t.TempDir()
	skillDir := filepath.Join(tmp, "skills", "lowercase")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := `---
name: lowercase-skill
description: Uses lowercase skill.md.
---
Body.
`
	if err := os.WriteFile(filepath.Join(skillDir, "skill.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	reg := NewRegistry(tmp)
	if err := reg.Discover(); err != nil {
		t.Fatal(err)
	}

	if _, ok := reg.Get("lowercase-skill"); !ok {
		t.Fatal("lowercase skill.md not discovered")
	}
}
