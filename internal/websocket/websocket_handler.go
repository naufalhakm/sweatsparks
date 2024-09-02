package websockets

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
	db  *sql.DB
}

func NewHandler(h *Hub, db *sql.DB) *Handler {
	return &Handler{
		hub: h,
		db:  db,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) ServeWs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	roomID := vars["roomID"]
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan *Message, 256),
		RoomID: roomID,
		Sender: userID,
	}
	client.Hub.Register <- client

	go client.WriteMessage()
	go client.ReadMessage(h.db)
}
