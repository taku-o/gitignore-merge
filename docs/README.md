# gitignore-merge

A command-line tool that intelligently merges multiple `.gitignore` files on a per-section basis.

[日本語](README.ja.md)

## Installation

```bash
go install github.com/taku-o/gitignore-merge/cmd/gitignore-merge@latest
```

Or build from source:

```bash
git clone <repository-url>
cd gitignore-merge
make build
```

## Usage

```bash
gitignore-merge <file1> <file2> [file3...]
```

- Specify two or more `.gitignore` file paths as arguments
- The merge result is printed to stdout
- Use redirection to write to a file

```bash
gitignore-merge project.gitignore template.gitignore > .gitignore
```

## Merge Rules

### Section Matching

Comment lines starting with `#` are recognized as section headers. Section names are compared by stripping all leading `#` characters and spaces from the first line of the header.

The following are all treated as the same section name "Node":

```
# Node
## Node
#Node
```

### Merge Priority

The first file serves as the base, and subsequent files are merged into it in order.

- **Same-name sections**: Patterns from subsequent files are added to the matching section in the base
- **New sections**: Sections not present in the base are appended to the end

### Deduplication

Exactly matching patterns within the same section are deduplicated.

### Conflict Resolution

When a pattern and its negation conflict (e.g., `path` and `!path`), the base file (first file) takes priority.

Examples:
- If the base has `*.log`, a subsequent `!*.log` is ignored
- If the base has `!*.log`, a subsequent `*.log` is ignored

## Example

### Input

**file1.gitignore** (base):

```gitignore
# Node
node_modules/
dist/

# Logs
*.log
!important.log

# OS
.DS_Store
Thumbs.db
```

**file2.gitignore**:

```gitignore
# Node
node_modules/
build/

# Logs
*.log
!debug.log

# IDE
.vscode/
.idea/
```

### Run

```bash
gitignore-merge file1.gitignore file2.gitignore
```

### Output

```gitignore
# Node
node_modules/
dist/

build/

# Logs
*.log
!important.log

!debug.log

# OS
.DS_Store
Thumbs.db

# IDE
.vscode/
.idea/
```

In this output:
- `# Node` section: `node_modules/` was deduplicated, `build/` was added
- `# Logs` section: `*.log` was deduplicated, `!debug.log` was added
- `# OS` section: exists only in the base, kept as-is
- `# IDE` section: exists only in file2, appended to the end
