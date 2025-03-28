package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDirectoryExists creates a directory if it doesn't exist
func EnsureDirectoryExists(path string) error {
	return os.MkdirAll(path, 0755)
}

// SaveJSONToFile saves data as JSON to the specified file
func SaveJSONToFile(path string, data interface{}) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := EnsureDirectoryExists(dir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal the data to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to the file
	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// SaveRawToFile saves raw data to the specified file
func SaveRawToFile(path string, data []byte) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := EnsureDirectoryExists(dir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to the file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// LoadJSONFromFile loads JSON data from a file into the provided destination
func LoadJSONFromFile(path string, dest interface{}) error {
	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal the JSON
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// FileExists checks if a file exists and is not a directory
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirectoryExists checks if a directory exists
func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
