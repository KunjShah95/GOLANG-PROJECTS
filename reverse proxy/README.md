# Go File Manager

A command-line file manager with basic file operations.

## Features
- List directory contents with metadata
- Create files/directories
- Delete files/directories
- Rename files/directories
- Preview file contents

## Installation
```bash
go build -o fm.exe
```

## Usage
```bash
# List directory contents
fm.exe list -path ./documents

# Create a new directory
fm.exe create -path ./new_folder -dir

# Create a new file
fm.exe create -path ./new_file.txt

# Delete a file/directory
fm.exe delete -path ./old_file.txt

# Rename a file
fm.exe rename -old ./file.txt -new ./renamed_file.txt

# View file contents
fm.exe view -path ./document.txt -lines 20
```
