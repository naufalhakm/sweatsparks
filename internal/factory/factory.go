package factory

import (
	"database/sql"
	"sweatsparks/internal/controllers"
	"sweatsparks/internal/repositories"
	"sweatsparks/internal/services"
)

type Provider struct {
	UserProvider    controllers.UserController
	MatchProvider   controllers.MatchController
	MessageProvider controllers.MessageController
	ProfileProvider controllers.ProfileController
	SwipeProvider   controllers.SwipeController
}

func InitFactory(db *sql.DB) *Provider {

	userRepo := repositories.NewUserRepository()
	userService := services.NeewUserService(db, userRepo)
	userController := controllers.NewUserController(userService)

	matchRepo := repositories.NewMatchRepository()
	matchService := services.NewMatchService(db, matchRepo)
	matchController := controllers.NewMatchController(matchService)

	messRepo := repositories.NewMessageRepository()
	messService := services.NewMessageService(db, messRepo)
	messController := controllers.NewMessageController(messService)

	profRepo := repositories.NewProfileRepository()
	profService := services.NewProfileService(db, profRepo)
	profController := controllers.NewProfileController(profService)

	swipeRepo := repositories.NewSwipeRepository()
	swipeService := services.NewSwipeService(db, swipeRepo)
	swipeController := controllers.NewSwipeController(swipeService)

	return &Provider{
		UserProvider:    userController,
		MatchProvider:   matchController,
		MessageProvider: messController,
		ProfileProvider: profController,
		SwipeProvider:   swipeController,
	}
}
