#!/usr/bin/env python3
"""
File Search Helper Script for WorkWise

This script provides file search functionality similar to macOS Spotlight,
allowing search by name, extension, content, size, and modification time.
"""

import os
import sys
import argparse
import fnmatch
import time
from pathlib import Path
from datetime import datetime, timedelta
import json


def human_readable_size(size):
    """Convert bytes to human-readable format."""
    for unit in ['B', 'KB', 'MB', 'GB', 'TB']:
        if size < 1024.0:
            return f"{size:.1f} {unit}"
        size /= 1024.0
    return f"{size:.1f} PB"


def human_readable_time(timestamp):
    """Convert timestamp to human-readable relative time."""
    now = datetime.now()
    file_time = datetime.fromtimestamp(timestamp)
    delta = now - file_time
    
    if delta.days > 365:
        years = delta.days // 365
        return f"{years} year{'s' if years > 1 else ''} ago"
    elif delta.days > 30:
        months = delta.days // 30
        return f"{months} month{'s' if months > 1 else ''} ago"
    elif delta.days > 0:
        return f"{delta.days} day{'s' if delta.days > 1 else ''} ago"
    elif delta.seconds > 3600:
        hours = delta.seconds // 3600
        return f"{hours} hour{'s' if hours > 1 else ''} ago"
    elif delta.seconds > 60:
        minutes = delta.seconds // 60
        return f"{minutes} minute{'s' if minutes > 1 else ''} ago"
    else:
        return "just now"


def matches_pattern(filename, pattern, case_sensitive=False):
    """Check if filename matches the given pattern."""
    if not case_sensitive:
        filename = filename.lower()
        pattern = pattern.lower()
    return fnmatch.fnmatch(filename, pattern)


def search_in_file(filepath, search_text, case_sensitive=False):
    """Search for text content in a file. Returns True if found."""
    try:
        with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
            content = f.read()
            if not case_sensitive:
                content = content.lower()
                search_text = search_text.lower()
            return search_text in content
    except (IOError, OSError, UnicodeDecodeError):
        # Skip files that can't be read or aren't text files
        return False


def search_files(
    search_path='.',
    name_pattern=None,
    extensions=None,
    content=None,
    max_depth=None,
    modified_days=None,
    min_size=None,
    max_size=None,
    max_results=100,
    case_sensitive=False
):
    """
    Search for files matching the given criteria.
    
    Args:
        search_path: Starting directory for search
        name_pattern: File name pattern (supports wildcards)
        extensions: List of file extensions to filter (e.g., ['.txt', '.pdf'])
        content: Text to search for within files
        max_depth: Maximum directory depth to search
        modified_days: Find files modified within N days
        min_size: Minimum file size in bytes
        max_size: Maximum file size in bytes
        max_results: Maximum number of results to return
        case_sensitive: Enable case-sensitive search
    
    Returns:
        List of matching file paths with metadata
    """
    results = []
    search_path = os.path.expanduser(search_path)
    
    if not os.path.exists(search_path):
        return {"error": f"Path does not exist: {search_path}"}
    
    # Calculate cutoff time for modified_days filter
    cutoff_time = None
    if modified_days is not None:
        cutoff_time = (datetime.now() - timedelta(days=modified_days)).timestamp()
    
    # Walk the directory tree
    for root, dirs, files in os.walk(search_path):
        # Check depth limit
        if max_depth is not None:
            depth = root[len(search_path):].count(os.sep)
            if depth >= max_depth:
                dirs.clear()  # Don't recurse deeper
        
        for filename in files:
            # Stop if we've reached max results
            if len(results) >= max_results:
                break
            
            filepath = os.path.join(root, filename)
            
            try:
                # Get file stats
                stat = os.stat(filepath)
                file_size = stat.st_size
                file_mtime = stat.st_mtime
                
                # Apply filters
                
                # Name pattern filter
                if name_pattern and not matches_pattern(filename, name_pattern, case_sensitive):
                    continue
                
                # Extension filter
                if extensions:
                    file_ext = os.path.splitext(filename)[1].lower()
                    if not any(file_ext == ext.lower() for ext in extensions):
                        continue
                
                # Size filters
                if min_size is not None and file_size < min_size:
                    continue
                if max_size is not None and file_size > max_size:
                    continue
                
                # Modified time filter
                if cutoff_time is not None and file_mtime < cutoff_time:
                    continue
                
                # Content search filter (more expensive, do last)
                if content and not search_in_file(filepath, content, case_sensitive):
                    continue
                
                # File matches all criteria
                results.append({
                    'path': filepath,
                    'name': filename,
                    'size': file_size,
                    'size_human': human_readable_size(file_size),
                    'modified': file_mtime,
                    'modified_human': human_readable_time(file_mtime),
                    'directory': root
                })
                
            except (OSError, PermissionError):
                # Skip files we can't access
                continue
        
        # Stop if we've reached max results
        if len(results) >= max_results:
            break
    
    return results


def format_results(results, max_results):
    """Format search results for display."""
    if isinstance(results, dict) and 'error' in results:
        return results['error']
    
    if not results:
        return "No files found matching the search criteria."
    
    output = []
    output.append(f"Found {len(results)} file{'s' if len(results) != 1 else ''}")
    
    if len(results) >= max_results:
        output.append(f"(showing first {max_results} results, use --max-results to see more)")
    
    output.append("")
    
    # Group results by directory
    by_directory = {}
    for result in results:
        directory = result['directory']
        if directory not in by_directory:
            by_directory[directory] = []
        by_directory[directory].append(result)
    
    # Sort directories
    for directory in sorted(by_directory.keys()):
        output.append(f"{directory}/")
        for result in sorted(by_directory[directory], key=lambda x: x['name']):
            output.append(f"  - {result['name']} ({result['size_human']}, modified {result['modified_human']})")
        output.append("")
    
    return "\n".join(output)


def main():
    parser = argparse.ArgumentParser(
        description='Search for files on the local filesystem (Spotlight-like functionality)'
    )
    
    parser.add_argument('--name', help='File name or pattern (supports wildcards * and ?)')
    parser.add_argument('--extension', action='append', help='File extension to filter (e.g., .txt, .pdf)')
    parser.add_argument('--path', default='.', help='Starting directory for search (default: current directory)')
    parser.add_argument('--content', help='Search for text content within files')
    parser.add_argument('--max-depth', type=int, help='Maximum directory depth to search')
    parser.add_argument('--modified-days', type=int, help='Find files modified within N days')
    parser.add_argument('--min-size', type=int, help='Minimum file size in bytes')
    parser.add_argument('--max-size', type=int, help='Maximum file size in bytes')
    parser.add_argument('--max-results', type=int, default=100, help='Maximum number of results (default: 100)')
    parser.add_argument('--case-sensitive', action='store_true', help='Enable case-sensitive search')
    parser.add_argument('--json', action='store_true', help='Output results in JSON format')
    
    args = parser.parse_args()
    
    # Perform search
    results = search_files(
        search_path=args.path,
        name_pattern=args.name,
        extensions=args.extension,
        content=args.content,
        max_depth=args.max_depth,
        modified_days=args.modified_days,
        min_size=args.min_size,
        max_size=args.max_size,
        max_results=args.max_results,
        case_sensitive=args.case_sensitive
    )
    
    # Output results
    if args.json:
        print(json.dumps(results, indent=2))
    else:
        print(format_results(results, args.max_results))


if __name__ == "__main__":
    main()
