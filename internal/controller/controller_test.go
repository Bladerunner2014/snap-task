package controller

import (
	"os"
	"testing"
	_ "github.com/mattn/go-sqlite3"
)

func TestNew(t *testing.T) {
	dbPath := "test.db"
	defer os.Remove(dbPath)

	controller, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create new SQLiteController: %v", err)
	}
	defer controller.Close()

	if controller.DB == nil {
		t.Error("DB is nil")
	}

	// Check if the table was created
	var tableName string
	err = controller.DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='responses'").Scan(&tableName)
	if err != nil {
		t.Errorf("Table 'responses' was not created: %v", err)
	}
	if tableName != "responses" {
		t.Errorf("Expected table name 'responses', got '%s'", tableName)
	}
}

func TestStoreResponse(t *testing.T) {
	dbPath := "test.db"
	defer os.Remove(dbPath)

	controller, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create new SQLiteController: %v", err)
	}
	defer controller.Close()

	url := "https://example.com"
	response := "Test response"

	err = controller.StoreResponse(url, response)
	if err != nil {
		t.Fatalf("Failed to store response: %v", err)
	}

	// Verify the stored data
	var storedURL, storedResponse string
	err = controller.DB.QueryRow("SELECT url, response FROM responses WHERE url = ?", url).Scan(&storedURL, &storedResponse)
	if err != nil {
		t.Fatalf("Failed to retrieve stored response: %v", err)
	}

	if storedURL != url {
		t.Errorf("Expected URL '%s', got '%s'", url, storedURL)
	}
	if storedResponse != response {
		t.Errorf("Expected response '%s', got '%s'", response, storedResponse)
	}
}

func TestClose(t *testing.T) {
	dbPath := "test.db"
	defer os.Remove(dbPath)

	controller, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create new SQLiteController: %v", err)
	}

	err = controller.Close()
	if err != nil {
		t.Errorf("Failed to close database: %v", err)
	}

	// Try to ping the closed database
	err = controller.DB.Ping()
	if err == nil {
		t.Error("Expected error when pinging closed database, got nil")
	}
}

func TestNewWithInvalidPath(t *testing.T) {
	dbPath := "/invalid/path/test.db"

	_, err := New(dbPath)
	if err == nil {
		t.Error("Expected error when creating SQLiteController with invalid path, got nil")
	}
}

func TestStoreResponseWithClosedDB(t *testing.T) {
	dbPath := "test.db"
	defer os.Remove(dbPath)

	controller, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create new SQLiteController: %v", err)
	}

	err = controller.Close()
	if err != nil {
		t.Fatalf("Failed to close database: %v", err)
	}

	url := "https://example.com"
	response := "Test response"

	err = controller.StoreResponse(url, response)
	if err == nil {
		t.Error("Expected error when storing response with closed database, got nil")
	}
}