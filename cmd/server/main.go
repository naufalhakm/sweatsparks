package main

import (
	"log"
	"net/http"
	"sweatsparks/internal/config"
	"sweatsparks/internal/factory"
	"sweatsparks/internal/routes"
	websockets "sweatsparks/internal/websocket"
	"sweatsparks/pkg/database"

	"github.com/gorilla/mux"
)

func main() {
	config.LoadConfig()
	psqlDB, err := database.NewPqSQLClient()
	if err != nil {
		log.Fatal("Could not connect to MySQL:", err)
	}

	router := mux.NewRouter()
	hub := websockets.NewHub()
	go hub.Run()

	provider := factory.InitFactory(psqlDB)
	routes.RegisterRoutes(psqlDB, router, hub, provider)

	log.Printf("Server running on :%s\n", config.ENV.ServerPort)
	log.Fatal(http.ListenAndServe(":"+config.ENV.ServerPort, router))

}
