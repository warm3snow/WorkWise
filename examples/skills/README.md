# WorkWise Skills Examples

This directory contains example skills demonstrating the Anthropic Agent Skills Spec implementation in WorkWise.

## What are Skills?

Skills are folders of instructions, scripts, and resources that AI agents can discover and load dynamically to improve performance on specialized tasks. Skills teach the agent how to complete specific tasks in a repeatable way.

## Skill Structure

Each skill is a directory containing:
- `SKILL.md` (required) - YAML frontmatter + markdown instructions
- `scripts/` (optional) - Executable helper scripts
- `resources/` (optional) - Additional files and templates

## Using Skills in WorkWise

### 1. Enable Skills

Edit your `~/.workwise/config.yaml`:

```yaml
extensions:
  skills_enabled: true
  skills_paths:
    - "~/.workwise/skills"
    - "./examples/skills"
```

### 2. Create or Install Skills

Place skill directories in one of the configured paths. Each skill directory must contain a `SKILL.md` file.

### 3. Using Skills

When you interact with WorkWise, it will automatically:
1. Load all skills from configured paths
2. Match user requests with appropriate skills based on descriptions
3. Apply skill instructions to improve response quality

## Creating Custom Skills

See the `example-skill/` directory for a template. The minimum required structure is:

```
my-skill/
└── SKILL.md
```

With content:

```markdown
---
name: my-skill
description: Clear description of what this skill does and when to use it
---

# My Skill

[Add your instructions here that the agent will follow when this skill is active]

## Examples
- Example usage 1
- Example usage 2

## Guidelines
- Guideline 1
- Guideline 2
```

## Example Skills in This Directory

### example-skill
A template and reference implementation showing the proper structure and format for creating skills.

## Resources

- [Anthropic Skills Repository](https://github.com/anthropics/skills) - Official skills examples
- [Agent Skills Spec](https://github.com/anthropics/skills/blob/main/agent_skills_spec.md) - Technical specification
- [Creating Custom Skills](https://support.claude.com/en/articles/12512198-creating-custom-skills) - Documentation

## Skill Development Guidelines

1. **Be Concise** - The context window is shared; only include essential information
2. **Clear Descriptions** - Help the agent know when to use the skill
3. **Provide Examples** - Show concrete usage patterns
4. **Test Thoroughly** - Validate skills work as expected before deployment

## Notes

- Skill names must be lowercase alphanumeric with hyphens (e.g., `my-skill`)
- Skill name must match the directory name
- Both `name` and `description` are required in the YAML frontmatter
- Instructions are in markdown format and loaded into the agent's context when the skill is used
