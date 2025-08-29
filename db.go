package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"
        "math"

	_ "modernc.org/sqlite"
)

type Message struct {
    ID          int
    Name        string
    Content     string
    CreatedAt   time.Time
}

var db *sql.DB

func initDB() {
    // try to open sqlite db
    var err error
    db, err = sql.Open("sqlite", "./data/guestbook.db")
    if err != nil {
        log.Fatal(err)
    }

    // try to create table
    if _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id          INTEGER PRIMARY KEY AUTOINCREMENT,
            name        TEXT NOT NULL,
            content     TEXT NOT NULL,
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
    totalPages := int(math.Ceil(float64(totalRows) / 15.0))

    return totalPages, nil
}

func getMessages(page int) ([]Message, error) {
    var offset int = page * 15 // first page is 0

    rows, err := db.Query(`
        SELECT id, name, content, created_at
        FROM messages
        ORDER BY created_at DESC
        LIMIT 15
        OFFSET ` + strconv.Itoa(offset))
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

func addMessage(name string, content string) error {
    _, err := db.Exec(`INSERT INTO messages (name, content) VALUES (?, ?)`, name, content)
    return err
}
