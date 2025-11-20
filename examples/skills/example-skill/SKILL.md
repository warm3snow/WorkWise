---
name: example-skill
description: An example skill demonstrating the Anthropic Skills specification. Use this as a template when creating new skills for WorkWise.
license: MIT
metadata:
  version: "1.0"
  author: "WorkWise"
---

# Example Skill

This is an example skill that demonstrates the structure and format of Anthropic-style skills.

## Overview

Skills are folders containing instructions, scripts, and resources that AI agents can discover and load dynamically to perform better at specific tasks. Each skill must contain a `SKILL.md` file with YAML frontmatter and markdown instructions.

## When to Use This Skill

Use this skill as a reference when creating new skills for WorkWise. It demonstrates:
- Proper YAML frontmatter structure
- Markdown instruction formatting
- Optional script integration
- Resource organization

## Skill Structure

```
skill-name/
├── SKILL.md          # Required: Frontmatter + instructions
├── scripts/          # Optional: Executable scripts
│   └── helper.py
└── resources/        # Optional: Additional files
    └── template.txt
```

## Creating Instructions

When writing skill instructions:

1. **Be concise** - The context window is shared with other content
2. **Provide clear workflows** - Step-by-step procedures for the task
3. **Include examples** - Show expected inputs and outputs
4. **Set appropriate freedom** - Balance specificity with flexibility

## Examples

### Example 1: Basic Task
```
User: "Help me with [task]"
Agent: [Follows instructions from this skill]
```

### Example 2: Using Scripts
When a script is needed, reference it by path:
```python
# scripts/helper.py can be executed for complex operations
```

## Guidelines

- Instructions should be clear and actionable
- Focus on the "how" more than the "why"
- Consider edge cases and error handling
- Test skills thoroughly before deploying

## References

For more information, see:
- [Anthropic Skills Repository](https://github.com/anthropics/skills)
- [Agent Skills Spec](https://github.com/anthropics/skills/blob/main/agent_skills_spec.md)
- [Creating Custom Skills](https://support.claude.com/en/articles/12512198-creating-custom-skills)
