package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSkill(t *testing.T) {
	validSkill := `---
name: test-skill
description: A test skill for unit testing
license: MIT
metadata:
  version: "1.0"
---

# Test Skill

This is a test skill with instructions.

## Usage

Use this skill for testing purposes.
`

	skill, err := ParseSkill(validSkill)
	if err != nil {
		t.Fatalf("Failed to parse valid skill: %v", err)
	}

	if skill.Name != "test-skill" {
		t.Errorf("Expected name 'test-skill', got '%s'", skill.Name)
	}

	if skill.Description != "A test skill for unit testing" {
		t.Errorf("Expected description 'A test skill for unit testing', got '%s'", skill.Description)
	}

	if skill.License != "MIT" {
		t.Errorf("Expected license 'MIT', got '%s'", skill.License)
	}

	if skill.Instructions == "" {
		t.Error("Expected non-empty instructions")
	}
}

func TestParseSkillInvalidFormat(t *testing.T) {
	invalidSkill := `
# Missing frontmatter
This should fail.
`

	_, err := ParseSkill(invalidSkill)
	if err == nil {
		t.Error("Expected error for invalid skill format")
	}
}

func TestParseSkillMissingName(t *testing.T) {
	missingName := `---
description: A test skill without a name
---

# Test Skill
`

	_, err := ParseSkill(missingName)
	if err == nil {
		t.Error("Expected error for missing name")
	}
}

func TestParseSkillMissingDescription(t *testing.T) {
	missingDesc := `---
name: test-skill
---

# Test Skill
`

	_, err := ParseSkill(missingDesc)
	if err == nil {
		t.Error("Expected error for missing description")
	}
}

func TestIsValidSkillName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"test-skill", true},
		{"my-skill-123", true},
		{"skill", true},
		{"Test-Skill", false}, // Uppercase not allowed
		{"test_skill", false}, // Underscore not allowed
		{"test.skill", false}, // Period not allowed
		{"", false},           // Empty not allowed
	}

	for _, tt := range tests {
		result := isValidSkillName(tt.name)
		if result != tt.valid {
			t.Errorf("isValidSkillName(%q) = %v, want %v", tt.name, result, tt.valid)
		}
	}
}

func TestLoadSkill(t *testing.T) {
	// Create a temporary skill directory
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "test-skill")
	err := os.Mkdir(skillDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create skill directory: %v", err)
	}

	// Create SKILL.md
	skillContent := `---
name: test-skill
description: A test skill
---

# Test Skill

Instructions for the test skill.
`
	skillPath := filepath.Join(skillDir, "SKILL.md")
	err = os.WriteFile(skillPath, []byte(skillContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write SKILL.md: %v", err)
	}

	// Load the skill
	skill, err := LoadSkill(skillDir)
	if err != nil {
		t.Fatalf("Failed to load skill: %v", err)
	}

	if skill.Name != "test-skill" {
		t.Errorf("Expected name 'test-skill', got '%s'", skill.Name)
	}

	if skill.SkillPath != skillDir {
		t.Errorf("Expected SkillPath '%s', got '%s'", skillDir, skill.SkillPath)
	}
}

func TestLoadSkillNameMismatch(t *testing.T) {
	// Create a temporary skill directory
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "wrong-name")
	err := os.Mkdir(skillDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create skill directory: %v", err)
	}

	// Create SKILL.md with different name
	skillContent := `---
name: correct-name
description: A test skill
---

# Test Skill
`
	skillPath := filepath.Join(skillDir, "SKILL.md")
	err = os.WriteFile(skillPath, []byte(skillContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write SKILL.md: %v", err)
	}

	// Try to load the skill - should fail due to name mismatch
	_, err = LoadSkill(skillDir)
	if err == nil {
		t.Error("Expected error for name mismatch")
	}
}

func TestLoader(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()

	// Create two test skills
	for i, name := range []string{"skill-one", "skill-two"} {
		skillDir := filepath.Join(tmpDir, name)
		err := os.Mkdir(skillDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create skill directory: %v", err)
		}

		skillContent := `---
name: ` + name + `
description: Test skill ` + string(rune('A'+i)) + `
---

# Test Skill
`
		skillPath := filepath.Join(skillDir, "SKILL.md")
		err = os.WriteFile(skillPath, []byte(skillContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write SKILL.md: %v", err)
		}
	}

	// Create loader and load skills
	loader := NewLoader([]string{tmpDir})
	err := loader.LoadAll()
	if err != nil {
		t.Fatalf("Failed to load skills: %v", err)
	}

	// Check that both skills were loaded
	names := loader.List()
	if len(names) != 2 {
		t.Errorf("Expected 2 skills, got %d", len(names))
	}

	// Check that we can get a specific skill
	skill, err := loader.Get("skill-one")
	if err != nil {
		t.Errorf("Failed to get skill-one: %v", err)
	}
	if skill.Name != "skill-one" {
		t.Errorf("Expected name 'skill-one', got '%s'", skill.Name)
	}
}
