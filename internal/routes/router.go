package routes

import (
	"database/sql"
	"net/http"
	"sweatsparks/internal/controllers"
	"sweatsparks/internal/middleware"
	websockets "sweatsparks/internal/websocket"

	"github.com/gorilla/mux"
)

func RegisterRoutes(db *sql.DB, router *mux.Router, hub *websockets.Hub, userController controllers.UserController, matchController controllers.MatchController) {
	router.HandleFunc("/api/register", userController.Register).Methods("POST")
	router.HandleFunc("/api/login", userController.Login).Methods("POST")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/users", userController.GetAllUsers).Methods("GET")
	protected.HandleFunc("/match", matchController.CreateMatch).Methods("POST")
	protected.HandleFunc("/match", matchController.GetAllMatchUser).Methods("GET")
	protected.HandleFunc("/match/{userID}", matchController.GetDetailMatchUser).Methods("GET")

	wsHandler := websockets.NewHandler(hub, db)
	router.HandleFunc("/ws/{userID}/room/{roomID}", func(w http.ResponseWriter, r *http.Request) {
		wsHandler.ServeWs(w, r)
	})
}
