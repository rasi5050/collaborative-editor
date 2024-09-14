package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
)

var db *sql.DB

func initDB() {
    var err error
    connStr := "user=username dbname=mydb password=mypassword sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to the database")
}

// SaveDocument persists a document to the database
func saveDocument(docID string, content string) error {
    query := `INSERT INTO documents (id, content) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET content = $2`
    _, err := db.Exec(query, docID, content)
    return err
}

// LoadDocument retrieves a document from the database
func loadDocument(docID string) (*Document, error) {
    var content string
    query := `SELECT content FROM documents WHERE id=$1`
    err := db.QueryRow(query, docID).Scan(&content)
    if err != nil {
        return nil, err
    }
    return &Document{Content: content, Version: 1}, nil
}

// SaveUser persists user credentials to the database
func saveUser(username, password string) error {
    hashedPassword, err := hashPassword(password)
    if err != nil {
        return err
    }
    query := `INSERT INTO users (username, password) VALUES ($1, $2)`
    _, err = db.Exec(query, username, hashedPassword)
    return err
}

// LoadUser retrieves user credentials from the database
func loadUser(username string) (string, error) {
    var passwordHash string
    query := `SELECT password FROM users WHERE username=$1`
    err := db.QueryRow(query, username).Scan(&passwordHash)
    return passwordHash, err
}
