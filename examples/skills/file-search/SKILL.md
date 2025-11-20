---
name: file-search
description: Search for files and directories on the local filesystem by name, path, content, or other attributes. Similar to macOS Spotlight functionality. Use this skill when users ask to find, locate, or search for files on their computer.
license: MIT
metadata:
  version: "1.0"
  author: "WorkWise"
---

# File Search Skill

This skill enables searching for files and directories on the local filesystem, providing functionality similar to macOS Spotlight or Windows Search.

## When to Use This Skill

Use this skill when the user wants to:
- Find files by name or pattern (e.g., "find all PDF files")
- Locate files by extension (e.g., "search for .txt files")
- Search for files in specific directories
- Find recently modified files
- Search for files containing specific text patterns
- Discover files by size or other attributes

## Capabilities

### Name-based Search
Search for files and directories by name or pattern:
- Exact matches: Find files with exact names
- Pattern matches: Use wildcards (* and ?) for flexible matching
- Case-insensitive searches by default

### Path-based Search
Search within specific directories:
- Recursive search through subdirectories
- Limit search depth if needed
- Search from specific starting points

### Extension-based Search
Find files by their extensions:
- Single extension (e.g., .pdf, .txt, .jpg)
- Multiple extensions (e.g., .go, .py, .js)

### Content Search
Search for text within files (when appropriate):
- Full-text search in text-based files
- Case-sensitive or case-insensitive matching

### Attribute-based Search
Find files by attributes:
- Size (files larger/smaller than specified size)
- Modified time (recently modified files)
- File type (regular files, directories, symlinks)

## Usage Instructions

### Basic File Search by Name

When a user asks to find files by name:
1. Determine the search criteria (exact name or pattern)
2. Ask for the starting directory if not specified (default to home directory or current directory)
3. Use the `search_files.py` script with appropriate parameters
4. Present results in a clear, organized format

Example command:
```bash
python3 scripts/search_files.py --name "report*.pdf" --path ~/Documents
```

### Search by Extension

When a user wants to find files by extension:
1. Identify the extension(s) to search for
2. Determine the search path
3. Execute search with extension filter

Example command:
```bash
python3 scripts/search_files.py --extension ".txt" --path ~/
```

### Search by Content

When searching for text within files:
1. Get the text pattern to search for
2. Determine which file types to search in
3. Specify the search directory
4. Execute content search

Example command:
```bash
python3 scripts/search_files.py --content "TODO" --extension ".py" --path ~/projects
```

### Recently Modified Files

When user wants to find recent files:
1. Determine the time frame (e.g., last 7 days)
2. Specify the directory to search
3. Execute with time filter

Example command:
```bash
python3 scripts/search_files.py --modified-days 7 --path ~/Documents
```

## Script Reference

The skill includes a helper script `scripts/search_files.py` with the following options:

```
--name          : File name or pattern (supports wildcards * and ?)
--extension     : File extension to filter (e.g., .txt, .pdf)
--path          : Starting directory for search (default: current directory)
--content       : Search for text content within files
--max-depth     : Maximum directory depth to search (default: unlimited)
--modified-days : Find files modified within N days
--min-size      : Minimum file size in bytes
--max-size      : Maximum file size in bytes
--max-results   : Maximum number of results to return (default: 100)
--case-sensitive: Enable case-sensitive search
```

## Response Format

When presenting search results:
1. **Summary**: Total number of matches found
2. **Results**: List of matching files with:
   - Full path
   - File size (in human-readable format)
   - Last modified date
   - File type (if relevant)
3. **Organization**: Group by directory or file type if helpful
4. **Limitations**: If max results reached, inform the user

Example response format:
```
Found 15 files matching "report*.pdf":

Documents/
  - report_2024.pdf (2.3 MB, modified 2 days ago)
  - report_draft.pdf (1.1 MB, modified 1 week ago)

Work/
  - annual_report.pdf (5.6 MB, modified 1 month ago)
  ...

Showing first 100 results. Use --max-results to see more.
```

## Best Practices

1. **Ask for clarification** if search criteria are ambiguous
2. **Set reasonable defaults**: 
   - Max 100 results unless user specifies more
   - Recursive search by default
   - Case-insensitive by default
3. **Provide helpful suggestions** if no results found:
   - Check spelling
   - Try broader search patterns
   - Search in different directories
4. **Security considerations**:
   - Don't search system directories without explicit user request
   - Respect file permissions
   - Warn about large searches that might take time
5. **Performance tips**:
   - Suggest more specific search criteria for broad searches
   - Offer to limit search depth for faster results

## Example Interactions

### Example 1: Find all Python files in a project
```
User: "Find all Python files in my projects folder"
Agent: Searching for Python files in ~/projects...
[Executes: python3 scripts/search_files.py --extension ".py" --path ~/projects]
[Returns organized list of .py files]
```

### Example 2: Find recent documents
```
User: "Show me documents I modified in the last week"
Agent: Searching for recently modified files...
[Executes: python3 scripts/search_files.py --modified-days 7 --path ~/Documents]
[Returns list with modification dates]
```

### Example 3: Find files containing specific text
```
User: "Find all files containing 'TODO' in my code"
Agent: Searching for 'TODO' in code files...
[Executes: python3 scripts/search_files.py --content "TODO" --extension ".py" ".go" ".js" --path ~/projects]
[Returns files with matches]
```

### Example 4: Find large files
```
User: "Find files larger than 100MB in my downloads"
Agent: Searching for large files...
[Executes: python3 scripts/search_files.py --min-size 104857600 --path ~/Downloads]
[Returns list sorted by size]
```

## Limitations

- Content search is limited to text-based files
- Very broad searches may take time and should be warned about
- Maximum results default prevents overwhelming output
- Some system files may not be accessible due to permissions

## Error Handling

If search fails:
1. Check if the search path exists and is accessible
2. Verify the search pattern is valid
3. Ensure the script has necessary permissions
4. Provide helpful error messages to the user
5. Suggest alternatives or corrections
