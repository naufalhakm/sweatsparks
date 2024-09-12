package routes

import (
	"database/sql"
	"net/http"
	"sweatsparks/internal/factory"
	"sweatsparks/internal/middleware"
	websockets "sweatsparks/internal/websocket"

	"github.com/gorilla/mux"
)

func RegisterRoutes(db *sql.DB, router *mux.Router, hub *websockets.Hub, provider *factory.Provider) {
	router.HandleFunc("/api/auth/register", provider.UserProvider.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", provider.UserProvider.Login).Methods("POST")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/users", provider.UserProvider.GetAllUsers).Methods("GET")

	protected.HandleFunc("/matches", provider.MatchProvider.CreateMatch).Methods("POST")
	protected.HandleFunc("/matches", provider.MatchProvider.GetAllMatchUser).Methods("GET")
	protected.HandleFunc("/matches/{userID}", provider.MatchProvider.GetDetailMatchUser).Methods("GET")

	protected.HandleFunc("/profiles", provider.ProfileProvider.GetAllProfile).Methods("GET")
	protected.HandleFunc("/profiles", provider.ProfileProvider.CreateProfile).Methods("POST")
	protected.HandleFunc("/profiles/{userID}", provider.ProfileProvider.GetDetailProfile).Methods("GET")
	protected.HandleFunc("/profiles/{userID}", provider.ProfileProvider.UpdateProfile).Methods("PATCH")

	protected.HandleFunc("/messages/{matchID}", provider.MessageProvider.GetMessageByMatchID).Methods("GET")

	protected.HandleFunc("/swipes", provider.SwipeProvider.CreateSwipe).Methods("POST")
	protected.HandleFunc("/swipes/{swiperID}", provider.SwipeProvider.GetSwipeAll).Methods("GET")
	protected.HandleFunc("/swipes/{swiperID}/swipee/{swipeeID}", provider.SwipeProvider.GetSwipeDetail).Methods("GET")

	wsHandler := websockets.NewHandler(hub, db)
	router.HandleFunc("/ws/{userID}/room/{roomID}", func(w http.ResponseWriter, r *http.Request) {
		wsHandler.ServeWs(w, r)
	})
}
