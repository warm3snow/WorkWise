package skills

import (
	"context"
	"fmt"
)

// Skill represents a skill that can be executed by the agent
type Skill interface {
	// Name returns the skill name
	Name() string

	// Description returns the skill description
	Description() string

	// Parameters returns the skill parameters schema
	Parameters() map[string]interface{}

	// Execute executes the skill with given parameters
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

// Registry manages skills
type Registry struct {
	skills map[string]Skill
}

// NewRegistry creates a new skills registry
func NewRegistry() *Registry {
	return &Registry{
		skills: make(map[string]Skill),
	}
}

// Register registers a skill
func (r *Registry) Register(skill Skill) error {
	name := skill.Name()
	if _, exists := r.skills[name]; exists {
		return fmt.Errorf("skill %s already registered", name)
	}
	r.skills[name] = skill
	return nil
}

// Get retrieves a skill by name
func (r *Registry) Get(name string) (Skill, error) {
	skill, exists := r.skills[name]
	if !exists {
		return nil, fmt.Errorf("skill %s not found", name)
	}
	return skill, nil
}

// List returns all registered skill names
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.skills))
	for name := range r.skills {
		names = append(names, name)
	}
	return names
}

// Execute executes a skill by name
func (r *Registry) Execute(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	skill, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	return skill.Execute(ctx, params)
}

// BaseSkill provides a base implementation for skills
type BaseSkill struct {
	name        string
	description string
	parameters  map[string]interface{}
}

// NewBaseSkill creates a new base skill
func NewBaseSkill(name, description string, parameters map[string]interface{}) *BaseSkill {
	return &BaseSkill{
		name:        name,
		description: description,
		parameters:  parameters,
	}
}

// Name returns the skill name
func (s *BaseSkill) Name() string {
	return s.name
}

// Description returns the skill description
func (s *BaseSkill) Description() string {
	return s.description
}

// Parameters returns the skill parameters
func (s *BaseSkill) Parameters() map[string]interface{} {
	return s.parameters
}

// Note: This is a foundational structure for future skills support.
// Concrete skills will be implemented as needed.
