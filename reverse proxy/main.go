package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Available commands: list, create, delete, rename, view")
		return
	}

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listPath := listCmd.String("path", ".", "Directory path")

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createPath := createCmd.String("path", "", "File/directory path")
	isDir := createCmd.Bool("dir", false, "Create directory")

	delCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	delPath := delCmd.String("path", "", "File/directory to delete")

	renameCmd := flag.NewFlagSet("rename", flag.ExitOnError)
	oldName := renameCmd.String("old", "", "Original path")
	newName := renameCmd.String("new", "", "New path")

	viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
	viewPath := viewCmd.String("path", "", "File to view")
	lines := viewCmd.Int("lines", 10, "Number of lines to preview")

	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		listDirectory(*listPath)
	case "create":
		createCmd.Parse(os.Args[2:])
		createPathHandler(*createPath, *isDir)
	case "delete":
		delCmd.Parse(os.Args[2:])
		deleteHandler(*delPath)
	case "rename":
		renameCmd.Parse(os.Args[2:])
		renameHandler(*oldName, *newName)
	case "view":
		viewCmd.Parse(os.Args[2:])
		viewFile(*viewPath, *lines)
	default:
		fmt.Println("Available commands: list, create, delete, rename, view")
	}
}

func listDirectory(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error listing directory: %v\n", err)
		return
	}

	fmt.Printf("Contents of %s:\n", path)
	for _, entry := range entries {
		info, _ := entry.Info()
		fmt.Printf("%s\t%s\t%d bytes\n",
			entry.Name(),
			info.Mode().String(),
			info.Size())
	}
}

func createPathHandler(path string, isDir bool) {
	if path == "" {
		fmt.Println("Path is required")
		return
	}

	if isDir {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}
		fmt.Printf("Created directory: %s\n", path)
	} else {
		file, err := os.Create(path)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}
		defer file.Close()
		fmt.Printf("Created file: %s\n", path)
	}
}

func deleteHandler(path string) {
	if path == "" {
		fmt.Println("Path is required")
		return
	}

	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Error deleting: %v\n", err)
		return
	}
	fmt.Printf("Deleted: %s\n", path)
}

func renameHandler(oldPath, newPath string) {
	if oldPath == "" || newPath == "" {
		fmt.Println("Both old and new paths are required")
		return
	}

	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Printf("Error renaming: %v\n", err)
		return
	}
	fmt.Printf("Renamed '%s' to '%s'\n", oldPath, newPath)
}

func viewFile(path string, lines int) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	contentStr := string(content)
	fileLines := strings.Split(contentStr, "\n")

	fmt.Printf("First %d lines of %s:\n", lines, path)
	for i := 0; i < lines && i < len(fileLines); i++ {
		fmt.Println(fileLines[i])
	}
}
