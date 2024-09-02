package main

import (
	"log"
	"net/http"
	"sweatsparks/internal/config"
	"sweatsparks/internal/controllers"
	"sweatsparks/internal/repositories"
	"sweatsparks/internal/routes"
	"sweatsparks/internal/services"
	websockets "sweatsparks/internal/websocket"
	"sweatsparks/pkg/database"

	"github.com/gorilla/mux"
)

func main() {
	config.LoadConfig()
	mysqlDB, err := database.NewMySQLClient()
	if err != nil {
		log.Fatal("Could not connect to MySQL:", err)
	}

	userRepo := repositories.NewUserRepository()
	userService := services.NeewUserService(mysqlDB, userRepo)
	userController := controllers.NewUserController(userService)

	matchRepo := repositories.NewMatchRepository()
	matchService := services.NewMatchService(mysqlDB, matchRepo)
	matchController := controllers.NewMatchController(matchService)

	router := mux.NewRouter()
	hub := websockets.NewHub()
	go hub.Run()

	routes.RegisterRoutes(mysqlDB, router, hub, userController, matchController)

	log.Printf("Server running on :%s\n", config.ENV.ServerPort)
	log.Fatal(http.ListenAndServe(":"+config.ENV.ServerPort, router))

}
