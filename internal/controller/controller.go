package controller

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

type SQLiteController struct {
    DB *sql.DB
}

func New(dbPath string) (*SQLiteController, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS responses (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            url TEXT,
            response TEXT,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        return nil, err
    }

    return &SQLiteController{DB: db}, nil
}

func (c *SQLiteController) StoreResponse(url, response string) error {
    _, err := c.DB.Exec("INSERT INTO responses (url, response) VALUES (?, ?)", url, response)
    if err != nil {
        log.Printf("Error storing response in database: %v", err)
        return err
    }
    return nil
}

func (c *SQLiteController) Close() error {
    return c.DB.Close()
}