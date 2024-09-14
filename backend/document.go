package main

import (
    "io/ioutil"
    "os"
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

type Document struct {
    Content     string
    Version     int
    LastEdited  time.Time
    Clients     map[*websocket.Conn]bool
    EditMutex   sync.Mutex
    UserCount   int
}

// Single global document reference
var doc = &Document{
    Content:    "",
    Version:    1,
    LastEdited: time.Now(),
    Clients:    make(map[*websocket.Conn]bool),
    UserCount:  0,
}

// AddClient adds a new WebSocket client to the document
func (d *Document) AddClient(conn *websocket.Conn) {
    d.EditMutex.Lock()
    defer d.EditMutex.Unlock()

    d.Clients[conn] = true
    d.UserCount++
    d.BroadcastUserCount()  // Notify all clients about the updated user count
}

// RemoveClient removes a WebSocket client from the document
func (d *Document) RemoveClient(conn *websocket.Conn) {
    d.EditMutex.Lock()
    defer d.EditMutex.Unlock()

    delete(d.Clients, conn)
    d.UserCount--
    d.BroadcastUserCount()  // Notify all clients about the updated user count
}

// BroadcastUserCount sends the current user count to all connected clients
func (d *Document) BroadcastUserCount() {
    for client := range d.Clients {
        err := client.WriteJSON(map[string]interface{}{
            "userCount": d.UserCount,
        })
        if err != nil {
            client.Close()
            delete(d.Clients, client)
        }
    }
}

// UpdateDocument updates the document content and last edited time
func (d *Document) UpdateDocument(newContent string) {
    d.EditMutex.Lock()
    defer d.EditMutex.Unlock()

    d.Content = newContent
    d.Version++
    d.LastEdited = time.Now()  // Update last edited time
}

// BroadcastUpdate sends the latest document content and last edited time to all connected clients
func (d *Document) BroadcastUpdate() {
    d.EditMutex.Lock()
    defer d.EditMutex.Unlock()

    for client := range d.Clients {
        err := client.WriteJSON(map[string]interface{}{
            "content":    d.Content,
            "lastEdited": d.LastEdited.Format(time.RFC3339),  // Send the timestamp in a standard format
        })
        if err != nil {
            client.Close()
            delete(d.Clients, client)
        }
    }
}

// SaveDocumentToFile saves the document content to a file
func saveDocumentToFile(content string) error {
    filename := "documents/single-doc.txt"
    return ioutil.WriteFile(filename, []byte(content), 0644)
}

// LoadDocumentFromFile loads the document content from a file
func loadDocumentFromFile() (string, error) {
    filename := "documents/single-doc.txt"
    content, err := ioutil.ReadFile(filename)
    if os.IsNotExist(err) {
        return "", nil
    }
    return string(content), err
}
