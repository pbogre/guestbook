package main

import (
	"database/sql"
	"log"
	"time"
        "math"

	_ "modernc.org/sqlite"
)

type Message struct {
    ID          int
    Name        string
    Content     string
    RemoteAddr  string
    CreatedAt   time.Time
}

var db *sql.DB

func initDB(path string) {
    // try to open sqlite db
    var err error
    db, err = sql.Open("sqlite", path)
    if err != nil {
        log.Fatal(err)
    }

    // try to create table
    if _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id          INTEGER PRIMARY KEY AUTOINCREMENT,
            name        TEXT NOT NULL,
            content     TEXT NOT NULL,
            remote_addr TEXT NOT NULL UNIQUE,
            created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
        );
    `); err != nil {
        log.Fatal(err)
    }
}

func getTotalPages() (int, error) {
    var totalRows int
    err := db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&totalRows)
    if err != nil {
        return 0, err
    }

    totalPages := math.Ceil(float64(totalRows) / float64(GuestbookConfig.EntriesPerPage))
    return int(math.Max(totalPages, 1)), nil // totalPages is never 0
}

func canRemoteAddrWrite(remoteAddr string) (bool, error) {
    var exists int
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM messages WHERE remote_addr = ?)", remoteAddr).Scan(&exists)
    if err != nil {
        log.Print(err)
        return false, err
    }

    return exists == 0, nil
}

func getMessages(page int) ([]Message, error) {
    entriesPerPage := GuestbookConfig.EntriesPerPage
    var offset int = page * entriesPerPage // first page is 0

    rows, err := db.Query(`
        SELECT id, name, content, created_at
        FROM messages
        ORDER BY created_at DESC
        LIMIT ?
        OFFSET ?`, entriesPerPage, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []Message
    for rows.Next() {
        var m Message
        if err := rows.Scan(&m.ID, &m.Name, &m.Content, &m.CreatedAt); err != nil {
            return nil, err
        }
        messages = append(messages, m)
    }

    return messages, nil
}

func postMessage(name, content, remote_addr string) error {
    _, err := db.Exec(`INSERT INTO messages (name, content, remote_addr) VALUES (?, ?, ?)`, name, content, remote_addr)
    return err
}
