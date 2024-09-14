package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer ws.Close()

    doc.AddClient(ws)
    defer doc.RemoveClient(ws)

    // Send the initial document state to the client, including the last edited time and user count
    ws.WriteJSON(map[string]interface{}{
        "content":    doc.Content,
        "lastEdited": doc.LastEdited.Format(time.RFC3339),
        "userCount":  doc.UserCount,
    })

    for {
        var msg map[string]string
        err := ws.ReadJSON(&msg)
        if err != nil {
            log.Printf("error: %v", err)
            break
        }

        newContent := msg["content"]
        doc.UpdateDocument(newContent)
        saveDocumentToFile(newContent)
        doc.BroadcastUpdate() // Broadcast the update to all clients
    }
}
