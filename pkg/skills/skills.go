package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Skill represents an Anthropic-style skill loaded from a SKILL.md file
// Skills are folders containing instructions, scripts, and resources that
// agents can discover and load dynamically to perform better at specific tasks.
//
// Reference: https://github.com/anthropics/skills
type Skill struct {
	// Metadata from YAML frontmatter
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	License      string            `yaml:"license,omitempty"`
	AllowedTools []string          `yaml:"allowed-tools,omitempty"`
	Metadata     map[string]string `yaml:"metadata,omitempty"`

	// Content from markdown body
	Instructions string

	// Path information
	SkillPath string // Path to the skill directory
}

// Loader manages loading and discovering skills from the filesystem
type Loader struct {
	skillPaths []string          // Directories to search for skills
	skills     map[string]*Skill // Loaded skills indexed by name
}

// NewLoader creates a new skills loader with the given search paths
func NewLoader(paths []string) *Loader {
	return &Loader{
		skillPaths: paths,
		skills:     make(map[string]*Skill),
	}
}

// LoadAll discovers and loads all skills from the configured paths
func (l *Loader) LoadAll() error {
	for _, basePath := range l.skillPaths {
		// Expand home directory if needed
		if strings.HasPrefix(basePath, "~/") {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to expand home directory: %w", err)
			}
			basePath = filepath.Join(home, basePath[2:])
		}

		// Check if path exists
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			continue // Skip non-existent paths
		}

		// Find all SKILL.md files
		err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && info.Name() == "SKILL.md" {
				skillDir := filepath.Dir(path)
				skill, err := LoadSkill(skillDir)
				if err != nil {
					return fmt.Errorf("failed to load skill from %s: %w", skillDir, err)
				}

				// Register the skill
				if _, exists := l.skills[skill.Name]; exists {
					return fmt.Errorf("duplicate skill name: %s", skill.Name)
				}
				l.skills[skill.Name] = skill
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to walk path %s: %w", basePath, err)
		}
	}

	return nil
}

// Get retrieves a loaded skill by name
func (l *Loader) Get(name string) (*Skill, error) {
	skill, exists := l.skills[name]
	if !exists {
		return nil, fmt.Errorf("skill not found: %s", name)
	}
	return skill, nil
}

// List returns all loaded skill names
func (l *Loader) List() []string {
	names := make([]string, 0, len(l.skills))
	for name := range l.skills {
		names = append(names, name)
	}
	return names
}

// GetAll returns all loaded skills
func (l *Loader) GetAll() []*Skill {
	skills := make([]*Skill, 0, len(l.skills))
	for _, skill := range l.skills {
		skills = append(skills, skill)
	}
	return skills
}

// LoadSkill loads a single skill from a directory containing a SKILL.md file
func LoadSkill(skillDir string) (*Skill, error) {
	skillPath := filepath.Join(skillDir, "SKILL.md")

	// Read the SKILL.md file
	content, err := os.ReadFile(skillPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SKILL.md: %w", err)
	}

	// Parse the skill
	skill, err := ParseSkill(string(content))
	if err != nil {
		return nil, err
	}

	// Set the skill path
	skill.SkillPath = skillDir

	// Validate that skill name matches directory name
	dirName := filepath.Base(skillDir)
	if skill.Name != dirName {
		return nil, fmt.Errorf("skill name '%s' does not match directory name '%s'", skill.Name, dirName)
	}

	return skill, nil
}

// ParseSkill parses a SKILL.md file content into a Skill structure
func ParseSkill(content string) (*Skill, error) {
	// Split frontmatter from body
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid SKILL.md format: missing YAML frontmatter")
	}

	// Parse YAML frontmatter
	skill := &Skill{}
	if err := yaml.Unmarshal([]byte(parts[1]), skill); err != nil {
		return nil, fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	// Validate required fields
	if skill.Name == "" {
		return nil, fmt.Errorf("skill name is required")
	}
	if skill.Description == "" {
		return nil, fmt.Errorf("skill description is required")
	}

	// Validate name format (lowercase alphanumeric + hyphens)
	if !isValidSkillName(skill.Name) {
		return nil, fmt.Errorf("invalid skill name format: %s (must be lowercase alphanumeric with hyphens)", skill.Name)
	}

	// Store markdown body as instructions
	skill.Instructions = strings.TrimSpace(parts[2])

	return skill, nil
}

// isValidSkillName checks if a skill name follows the specification
// (lowercase Unicode alphanumeric + hyphen)
func isValidSkillName(name string) bool {
	if name == "" {
		return false
	}

	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}

	return true
}

// GetScriptPath returns the path to a script within the skill directory
func (s *Skill) GetScriptPath(scriptName string) string {
	return filepath.Join(s.SkillPath, "scripts", scriptName)
}

// HasScript checks if a skill has a specific script
func (s *Skill) HasScript(scriptName string) bool {
	scriptPath := s.GetScriptPath(scriptName)
	_, err := os.Stat(scriptPath)
	return err == nil
}

// GetResourcePath returns the path to a resource within the skill directory
func (s *Skill) GetResourcePath(resourceName string) string {
	return filepath.Join(s.SkillPath, resourceName)
}

// Note: This implementation follows the Anthropic Agent Skills Spec.
// Skills are folder-based with SKILL.md files containing YAML frontmatter
// and markdown instructions. The agent loads these dynamically to improve
// performance on specialized tasks.
//
// For more information, see: https://github.com/anthropics/skills
