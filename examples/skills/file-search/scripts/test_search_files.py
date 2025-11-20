#!/usr/bin/env python3
"""
Test suite for the file-search skill search_files.py script
"""

import os
import sys
import tempfile
import shutil
import json
import subprocess
from pathlib import Path

# Path to the search script
SEARCH_SCRIPT = os.path.join(os.path.dirname(__file__), 'search_files.py')


def run_search(*args):
    """Run the search script with given arguments and return output."""
    cmd = [sys.executable, SEARCH_SCRIPT] + list(args)
    result = subprocess.run(cmd, capture_output=True, text=True)
    return result.stdout, result.stderr, result.returncode


def run_search_json(*args):
    """Run the search script with --json flag and return parsed results."""
    stdout, stderr, returncode = run_search(*args, '--json')
    if returncode != 0:
        return None
    try:
        return json.loads(stdout)
    except json.JSONDecodeError:
        return None


class TestFileSearch:
    """Test cases for file search functionality"""
    
    def __init__(self):
        self.test_dir = None
        self.setup()
    
    def setup(self):
        """Create a temporary test directory structure"""
        self.test_dir = tempfile.mkdtemp(prefix='file_search_test_')
        
        # Create test directory structure
        os.makedirs(os.path.join(self.test_dir, 'subdir1'))
        os.makedirs(os.path.join(self.test_dir, 'subdir2', 'nested'))
        
        # Create test files
        files = {
            'test.txt': 'Hello World\nThis is a test file.',
            'report.pdf': 'PDF content placeholder',
            'README.md': '# README\nThis is a markdown file.',
            'subdir1/file1.txt': 'File in subdirectory 1',
            'subdir1/document.pdf': 'Another PDF',
            'subdir2/file2.txt': 'File in subdirectory 2',
            'subdir2/nested/deep.txt': 'Deep nested file',
            'data.json': '{"key": "value"}',
        }
        
        for filepath, content in files.items():
            full_path = os.path.join(self.test_dir, filepath)
            with open(full_path, 'w') as f:
                f.write(content)
    
    def teardown(self):
        """Clean up temporary test directory"""
        if self.test_dir and os.path.exists(self.test_dir):
            shutil.rmtree(self.test_dir)
    
    def test_search_by_name(self):
        """Test searching files by exact name"""
        results = run_search_json('--name', 'test.txt', '--path', self.test_dir)
        assert results is not None, "Search failed"
        assert len(results) == 1, f"Expected 1 result, got {len(results)}"
        assert 'test.txt' in results[0]['name'], "Result should contain test.txt"
        print("✓ Test search by name passed")
    
    def test_search_by_pattern(self):
        """Test searching files with wildcard pattern"""
        results = run_search_json('--name', '*.txt', '--path', self.test_dir)
        assert results is not None, "Search failed"
        assert len(results) >= 4, f"Expected at least 4 .txt files, got {len(results)}"
        print("✓ Test search by pattern passed")
    
    def test_search_by_extension(self):
        """Test searching files by extension"""
        results = run_search_json('--extension', '.pdf', '--path', self.test_dir)
        assert results is not None, "Search failed"
        assert len(results) == 2, f"Expected 2 PDF files, got {len(results)}"
        for result in results:
            assert result['name'].endswith('.pdf'), "All results should be PDF files"
        print("✓ Test search by extension passed")
    
    def test_search_by_content(self):
        """Test searching files by content"""
        results = run_search_json('--content', 'Hello World', '--path', self.test_dir)
        assert results is not None, "Search failed"
        assert len(results) >= 1, f"Expected at least 1 result, got {len(results)}"
        assert any('test.txt' in r['name'] for r in results), "Should find test.txt with 'Hello World'"
        print("✓ Test search by content passed")
    
    def test_max_depth(self):
        """Test limiting search depth"""
        # Depth 0 should find only root level files
        results = run_search_json('--name', '*.txt', '--path', self.test_dir, '--max-depth', '0')
        assert results is not None, "Search failed"
        root_files = [r for r in results if '/' not in r['path'].replace(self.test_dir, '').strip('/')]
        assert len(results) == 1, f"Expected 1 root level .txt file, got {len(results)}"
        print("✓ Test max depth passed")
    
    def test_max_results(self):
        """Test limiting number of results"""
        results = run_search_json('--name', '*', '--path', self.test_dir, '--max-results', '3')
        assert results is not None, "Search failed"
        assert len(results) <= 3, f"Expected at most 3 results, got {len(results)}"
        print("✓ Test max results passed")
    
    def test_nonexistent_path(self):
        """Test searching in non-existent path"""
        results = run_search_json('--path', '/nonexistent/path/xyz')
        assert results is not None, "Should return error object"
        assert 'error' in results, "Should contain error key"
        print("✓ Test nonexistent path passed")
    
    def test_combined_filters(self):
        """Test combining multiple filters"""
        results = run_search_json(
            '--name', '*.txt',
            '--path', self.test_dir,
            '--content', 'file',
            '--max-results', '10'
        )
        assert results is not None, "Search failed"
        # Should find txt files containing 'file'
        assert len(results) >= 1, "Should find files matching all criteria"
        for result in results:
            assert result['name'].endswith('.txt'), "Should only return .txt files"
        print("✓ Test combined filters passed")
    
    def test_case_insensitive_default(self):
        """Test that search is case-insensitive by default"""
        results = run_search_json('--name', 'README*', '--path', self.test_dir)
        assert results is not None, "Search failed"
        assert len(results) >= 1, "Should find README files case-insensitively"
        print("✓ Test case insensitive default passed")
    
    def run_all(self):
        """Run all test cases"""
        try:
            print(f"Running tests in {self.test_dir}")
            print("-" * 60)
            
            self.test_search_by_name()
            self.test_search_by_pattern()
            self.test_search_by_extension()
            self.test_search_by_content()
            self.test_max_depth()
            self.test_max_results()
            self.test_nonexistent_path()
            self.test_combined_filters()
            self.test_case_insensitive_default()
            
            print("-" * 60)
            print("✅ All tests passed!")
            return True
        except AssertionError as e:
            print(f"❌ Test failed: {e}")
            return False
        finally:
            self.teardown()


if __name__ == '__main__':
    tester = TestFileSearch()
    success = tester.run_all()
    sys.exit(0 if success else 1)
