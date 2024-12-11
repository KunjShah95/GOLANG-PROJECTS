package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// define categories for file extensions
var fileCategories = map[string][]string{
	"Images":    {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg"},
	"Documents": {".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"},
	"Music":     {".mp3", ".wav", ".ogg", ".flac", ".aac"},
	"Videos":    {".mp4", ".avi", ".mkv", ".mov", ".wmv"},
	"Archives":  {".zip", ".tar", ".gz", ".7z", ".rar"},
	"Programs":  {".exe", ".msi", ".bat", ".sh", ".py", ".js", ".css", ".html", ".json"},
	"Other":     {}, // for files that don't match any category
}

func main() {
	// Set up logging
	logFile, err := os.OpenFile("file_organizer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could not open log file:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Get the path to organize
	fmt.Print("Enter the path to organize (e.g., C:\\ or D:\\Folder): ")
	var dirPath string
	fmt.Scanln(&dirPath)

	// Check if the path exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Println("The specified path does not exist.")
		log.Println("Error: The specified path does not exist:", dirPath)
		return
	}

	// Confirm if the user wants to organize the root of a drive
	if strings.HasSuffix(dirPath, ":\\") {
		fmt.Println("Warning: Organizing the root of a drive may affect system files.")
		fmt.Print("Are you sure you want to continue? (yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if !strings.EqualFold(confirm, "yes") {
			fmt.Println("Operation cancelled.")
			return
		}
	}

	// Restrict to only C: and D: drives
	if !strings.HasPrefix(dirPath, "C:\\") && !strings.HasPrefix(dirPath, "D:\\") {
		fmt.Println("This program only supports organizing files in C: or D: drives.")
		return
	}

	// Ask user for specific categories to organize
	fmt.Println("Available categories: Images, Documents, Music, Videos, Archives, Programs, Other")
	fmt.Print("Enter categories to organize (comma-separated, or leave blank for all): ")
	var categoriesInput string
	fmt.Scanln(&categoriesInput)
	categoriesToOrganize := parseCategories(categoriesInput)

	// Read files from the directory and subdirectories
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error accessing path:", path, err)
			return nil // continue walking
		}
		if info.IsDir() {
			return nil // skip directories
		}

		fileExt := strings.ToLower(filepath.Ext(info.Name()))
		category := getCategory(fileExt)

		// Check if the category is in the user's specified categories
		if len(categoriesToOrganize) > 0 && !contains(categoriesToOrganize, category) {
			return nil // skip this file if not in specified categories
		}

		// Create category folder if it doesn't exist
		categoryPath := filepath.Join(filepath.Dir(path), category)
		if _, err := os.Stat(categoryPath); os.IsNotExist(err) {
			err := os.Mkdir(categoryPath, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating folder:", err)
				log.Println("Error creating folder:", err)
				return nil
			}
		}

		// Move the file to the category folder
		newPath := filepath.Join(categoryPath, info.Name())
		err = os.Rename(path, newPath)
		if err != nil {
			fmt.Println("Error moving file:", info.Name(), err)
			log.Println("Error moving file:", info.Name(), err)
		} else {
			fmt.Println("Moved:", info.Name(), "to", category)
			log.Println("Moved:", info.Name(), "to", category)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
		log.Println("Error walking the directory:", err)
		return
	}

	fmt.Println("File organization complete!")
}

// getCategory returns the category for a given file extension
func getCategory(ext string) string {
	for category, extensions := range fileCategories {
		for _, validExt := range extensions {
			if ext == validExt {
				return category
			}
		}
	}
	return "Other"
}

// parseCategories splits the user input into a slice of categories
func parseCategories(input string) []string {
	if input == "" {
		return []string{} // return empty slice for all categories
	}
	categories := strings.Split(input, ",")
	for i := range categories {
		categories[i] = strings.TrimSpace(strings.Title(strings.ToLower(categories[i])))
	}
	return categories
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
