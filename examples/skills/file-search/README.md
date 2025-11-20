# File Search Skill

A local file search skill for WorkWise that provides Spotlight-like functionality for searching files on the local filesystem.

## Overview

This skill enables the AI agent to help users find files by name, extension, content, size, and modification time. It's similar to macOS Spotlight or Windows Search, but accessible through the WorkWise AI assistant.

## Features

- **Name-based search**: Find files by exact name or pattern using wildcards
- **Extension filtering**: Search for specific file types (.pdf, .txt, .jpg, etc.)
- **Content search**: Find files containing specific text
- **Attribute filtering**: Filter by file size and modification time
- **Flexible search depth**: Control how deep into directory structures to search
- **Human-readable output**: File sizes and modification times in readable format

## Installation

This skill is included in the WorkWise `examples/skills/` directory. To use it:

1. Enable skills in your WorkWise configuration (`~/.workwise/config.yaml`):
   ```yaml
   extensions:
     skills_enabled: true
     skills_paths:
       - "./examples/skills"
   ```

2. The skill will be automatically loaded when WorkWise starts

## Usage Examples

### Through WorkWise AI Agent

Once the skill is loaded, you can ask the AI agent to search for files naturally:

```
You: Find all PDF files in my Documents folder
Agent: [Uses file-search skill to locate PDFs]

You: Show me Python files modified in the last week
Agent: [Searches for .py files with recent modifications]

You: Find files containing "TODO" in my project
Agent: [Performs content search in project directory]
```

### Direct Script Usage

You can also run the search script directly:

```bash
# Search for files by name pattern
python3 scripts/search_files.py --name "*.pdf" --path ~/Documents

# Search for files by extension
python3 scripts/search_files.py --extension ".py" --path ~/projects

# Search for content in files
python3 scripts/search_files.py --content "TODO" --extension ".py" --path ~/projects

# Find recently modified files
python3 scripts/search_files.py --modified-days 7 --path ~/Documents

# Find large files
python3 scripts/search_files.py --min-size 10485760 --path ~/Downloads

# Combine multiple filters
python3 scripts/search_files.py --name "report*" --extension ".pdf" --modified-days 30
```

## Script Options

```
--name          : File name or pattern (supports wildcards * and ?)
--extension     : File extension to filter (e.g., .txt, .pdf)
--path          : Starting directory for search (default: current directory)
--content       : Search for text content within files
--max-depth     : Maximum directory depth to search
--modified-days : Find files modified within N days
--min-size      : Minimum file size in bytes
--max-size      : Maximum file size in bytes
--max-results   : Maximum number of results (default: 100)
--case-sensitive: Enable case-sensitive search
--json          : Output results in JSON format
```

## Testing

Run the test suite to verify functionality:

```bash
python3 scripts/test_search_files.py
```

This will run comprehensive tests including:
- Name-based search
- Pattern matching with wildcards
- Extension filtering
- Content search
- Search depth limits
- Result limits
- Error handling

## Performance Considerations

- Content search is slower than name/extension search as it reads file contents
- Large directory trees may take time to search; consider using `--max-depth` to limit scope
- Use specific search criteria to get faster, more relevant results
- The default limit of 100 results prevents overwhelming output

## Limitations

- Content search only works with text-based files
- Binary files are skipped during content search
- Some system files may not be accessible due to permissions
- Very large files may be slow to search through

## Security

- The script respects file system permissions
- It does not modify any files, only reads them
- Users should be cautious when searching system directories
- The AI agent is instructed to avoid searching sensitive directories without explicit user request

## Integration with WorkWise

When loaded by WorkWise, the skill provides:
- Comprehensive instructions for the AI agent on when and how to use file search
- Guidance on presenting results in a user-friendly format
- Best practices for asking clarifying questions
- Error handling and user feedback patterns

## Development

The skill consists of:
- `SKILL.md`: Instructions for the AI agent
- `scripts/search_files.py`: Python script for performing searches
- `scripts/test_search_files.py`: Test suite
- `README.md`: This documentation

To modify or extend the skill:
1. Update `search_files.py` for new search capabilities
2. Update `SKILL.md` to instruct the AI on using new features
3. Add tests in `test_search_files.py`
4. Update this README with new documentation

## License

MIT License - See the LICENSE file in the repository root for details.
