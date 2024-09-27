package controller

import (
    "testing"

    _ "github.com/mattn/go-sqlite3"
)

// TestNewSQLiteController tests that a new SQLiteController can be created
func TestNewSQLiteController(t *testing.T) {
    // Create an in-memory SQLite DB for testing
    dbPath := ":memory:"

    ctrl, err := New(dbPath)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if ctrl.DB == nil {
        t.Fatalf("Expected a valid database connection, got nil")
    }

    // Cleanup: close the DB
    ctrl.Close()
}

// TestStoreResponse tests that StoreResponse inserts data correctly
func TestStoreResponse(t *testing.T) {
    dbPath := ":memory:"

    // Create new controller with in-memory SQLite DB
    ctrl, err := New(dbPath)
    if err != nil {
        t.Fatalf("Error initializing SQLiteController: %v", err)
    }
    defer ctrl.Close()

    // Create the responses table (since this was commented out in the original code)
    _, err = ctrl.DB.Exec(`
        CREATE TABLE responses (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            url TEXT,
            response TEXT,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        t.Fatalf("Error creating table: %v", err)
    }

    // Store a response
    url := "http://example.com"
    response := "Example Response"
    err = ctrl.StoreResponse(url, response)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    // Query the database to ensure the data was inserted
    var storedURL, storedResponse string
    row := ctrl.DB.QueryRow("SELECT url, response FROM responses WHERE url = ?", url)
    err = row.Scan(&storedURL, &storedResponse)
    if err != nil {
        t.Fatalf("Error retrieving stored data: %v", err)
    }

    if storedURL != url || storedResponse != response {
        t.Fatalf("Expected (%s, %s), got (%s, %s)", url, response, storedURL, storedResponse)
    }
}

// TestClose tests that the Close method works without errors
func TestClose(t *testing.T) {
    dbPath := ":memory:"

    // Create a new controller
    ctrl, err := New(dbPath)
    if err != nil {
        t.Fatalf("Error initializing SQLiteController: %v", err)
    }

    // Close the controller
    err = ctrl.Close()
    if err != nil {
        t.Fatalf("Expected no error closing the DB, got %v", err)
    }

    // Ensure that the DB is closed by attempting to query
    err = ctrl.DB.Ping()
    if err == nil {
        t.Fatalf("Expected error on pinging a closed DB, got none")
    }
}
