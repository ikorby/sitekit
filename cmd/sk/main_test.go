package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_CreatesProjectStructure(t *testing.T) {
	tempDir := t.TempDir()
	args := []string{"sk", "new", "testsite"}

	err := run(args, tempDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	targetDir := filepath.Join(tempDir, "testsite")

	// Проверяем, что создались нужные директории
	expectedDirs := []string{
		filepath.Join(targetDir, "cmd", "newsite"),
		filepath.Join(targetDir, "templates", "layouts"),
		filepath.Join(targetDir, "static", "css"),
	}

	for _, dir := range expectedDirs {
		info, err := os.Stat(dir)
		if os.IsNotExist(err) {
			t.Errorf("expected directory %s to be created", dir)
		} else if !info.IsDir() {
			t.Errorf("expected %s to be a directory", dir)
		}
	}

	// Проверяем, что создались нужные файлы
	expectedFiles := []string{
		filepath.Join(targetDir, "go.mod"),
		filepath.Join(targetDir, ".env"),
		filepath.Join(targetDir, "cmd", "newsite", "main.go"),
		filepath.Join(targetDir, "templates", "layouts", "base.html"),
	}

	for _, file := range expectedFiles {
		info, err := os.Stat(file)
		if os.IsNotExist(err) {
			t.Errorf("expected file %s to be created", file)
		} else if info.IsDir() {
			t.Errorf("expected %s to be a file", file)
		}
	}
}

func TestRun_ValidatesArguments(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name string
		args []string
	}{
		{"NoArgs", []string{"sk"}},
		{"WrongCommand", []string{"sk", "build", "testsite"}},
		{"MissingProjectName", []string{"sk", "new"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.args, tempDir)
			if err == nil {
				t.Errorf("expected error for args %v, got nil", tt.args)
			}
		})
	}
}
