package agentskills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoader(t *testing.T) {
	content := `---
name: test-skill
description: A test skill for testing.
allowed-tools: [bash, read_file]
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
	if len(s.AllowedTools) != 2 || s.AllowedTools[0] != "bash" {
		t.Errorf("expected allowed-tools [bash, read_file], got %v", s.AllowedTools)
	}
	if s.Instructions != "# Instructions\nUse this skill for testing purposes." {
		t.Errorf("expected instructions, got %s", s.Instructions)
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
