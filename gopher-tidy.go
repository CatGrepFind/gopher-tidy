package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// main drives the application flow
func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("--- Application Cleaner (macOS Edition) ---")
	fmt.Print("Enter the name of the application to clean up (e.g., Docker): ")

	appName, _ := reader.ReadString('\n')
	appName = strings.TrimSpace(appName)

	if appName == "" {
		fmt.Println("Application name cannot be empty. Exiting.")
		return
	}

	fmt.Printf("\n🔍 Searching for files related to '%s'...\n", appName)
	foundFiles, err := findAssociatedFiles(appName)
	if err != nil {
		fmt.Printf("Error during file search: %v\n", err)
		return
	}

	if len(foundFiles) == 0 {
		fmt.Println("✅ No associated files found in common locations.")
		return
	}

	// Start the interactive deletion menu
	handleDeletion(foundFiles, reader)
}

// findAssociatedFiles searches common macOS directories for files/folders matching the app name.
func findAssociatedFiles(appName string) ([]string, error) {
	var pathsToSearch []string
	var foundFiles []string

	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("could not get current user: %w", err)
	}
	homeDir := currentUser.HomeDir

	// List of common directories where apps leave files
	pathsToSearch = []string{
		filepath.Join(homeDir, "Library", "Application Support"),
		filepath.Join(homeDir, "Library", "Caches"),
		filepath.Join(homeDir, "Library", "Preferences"),
		filepath.Join(homeDir, "Library", "Logs"),
		filepath.Join(homeDir, "Library", "Saved Application State"),
		"/Library/Application Support", // System-level directories
		"/Library/Caches",
		"/Library/LaunchAgents",
		"/Library/LaunchDaemons",
	}

	// Use a map to avoid duplicate entries
	uniquePaths := make(map[string]struct{})

	searchKeyword := strings.ToLower(strings.ReplaceAll(appName, " ", ""))

	for _, path := range pathsToSearch {
		// The WalkDir function traverses the directory tree
		err := filepath.WalkDir(path, func(currentPath string, d os.DirEntry, err error) error {
			if err != nil {
				// Silently ignore permission errors during search phase
				if os.IsPermission(err) {
					return nil
				}
				return err
			}

			// We only check directories at the top level of our search paths
			// and files within them. This avoids overly broad matches deep in unrelated folders.
			if filepath.Dir(currentPath) == path {
				baseName := strings.ToLower(strings.ReplaceAll(d.Name(), " ", ""))
				if strings.Contains(baseName, searchKeyword) {
					uniquePaths[currentPath] = struct{}{}
				}
			}

			return nil
		})
		if err != nil {
			// Don't stop if one path is inaccessible, just log it and continue
			fmt.Printf("⚠️  Could not fully search '%s': %v\n", path, err)
		}
	}

	for path := range uniquePaths {
		foundFiles = append(foundFiles, path)
	}

	return foundFiles, nil
}

// handleDeletion displays the interactive menu for deleting files.
func handleDeletion(files []string, reader *bufio.Reader) {
	for {
		fmt.Println("\nFound potential leftover files and folders:")
		for i, file := range files {
			fmt.Printf("  [%d] %s\n", i+1, file)
		}
		fmt.Println("\nEnter numbers to delete (e.g., 1 3 4), 'all' to delete everything, or 'quit' to exit.")
		fmt.Print("> ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "quit") {
			fmt.Println("Exiting without changes.")
			break
		}

		var filesToDelete []string
		if strings.EqualFold(input, "all") {
			filesToDelete = files
		} else {
			indices := strings.Fields(input)
			for _, idxStr := range indices {
				idx, err := strconv.Atoi(idxStr)
				if err != nil || idx < 1 || idx > len(files) {
					fmt.Printf("Invalid selection: '%s'. Please enter a valid number.\n", idxStr)
					continue
				}
				filesToDelete = append(filesToDelete, files[idx-1])
			}
		}

		if len(filesToDelete) > 0 {
			fmt.Println("\n--- DELETION SUMMARY ---")
			for _, file := range filesToDelete {
				fmt.Printf("  - %s\n", file)
			}
			fmt.Print("Proceed with deleting these items? [y/N]: ")

			confirm, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(confirm)) == "y" {
				deleteFiles(filesToDelete)
				// After deletion, exit the loop.
				// A more advanced version might rescan and continue.
				break
			} else {
				fmt.Println("Deletion cancelled.")
			}
		}
	}
}

// deleteFiles attempts to remove the specified files and folders.
func deleteFiles(files []string) {
	fmt.Println("\n🚀 Starting deletion...")
	for _, file := range files {
		fmt.Printf("Attempting to delete: %s\n", file)
		// os.RemoveAll deletes a path and any children it contains.
		err := os.RemoveAll(file)
		if err != nil {
			// This is the crucial part for handling permissions.
			if os.IsPermission(err) {
				fmt.Printf("❌ PERMISSION DENIED for: %s\n", file)
				fmt.Printf("   This file requires administrator privileges. To delete it, run this command in your terminal:\n")
				// We print the command for the user to run securely.
				fmt.Printf("   sudo rm -rf \"%s\"\n\n", file)
			} else {
				fmt.Printf("❌ ERROR deleting %s: %v\n", file, err)
			}
		} else {
			fmt.Printf("✅ DELETED: %s\n", file)
		}
	}
	fmt.Println("\n✨ Cleanup process finished.")
}
