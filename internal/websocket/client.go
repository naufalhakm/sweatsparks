package websockets

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan *Message
	RoomID string
	Sender string
}

type Message struct {
	RoomID  string `json:"room_id"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

type IncomingMessage struct {
	Content string `json:"content"`
	File    string `json:"file"`
}

func (c *Client) WriteMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			}

			c.Conn.WriteJSON(message)
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadMessage(db *sql.DB) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var incomingMsg IncomingMessage
		err = json.Unmarshal(message, &incomingMsg)
		if err != nil {
			log.Println("Error decoding JSON message:", err)
			continue
		}
		msg := &Message{
			RoomID:  c.RoomID,
			Sender:  c.Sender,
			Content: incomingMsg.Content,
			Time:    time.Now().Format(time.RFC3339),
		}

		if err := storeMessage(db, msg); err != nil {
			log.Printf("errror storing message %v", err)
		}

		c.Hub.Broadcast <- msg
	}
}

func storeMessage(db *sql.DB, msg *Message) error {
	query := `INSERT INTO messages (match_id,sender_id,content,sent_at) VALUES (?,?,?,?)`
	_, err := db.Exec(query, msg.RoomID, msg.Sender, msg.Content, msg.Time)
	return err
}
