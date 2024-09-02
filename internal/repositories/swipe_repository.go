package repositories

import (
	"context"
	"database/sql"
	"sweatsparks/internal/models"
)

type SwipeRepository interface {
	CreateSwipe(ctx context.Context, tx *sql.DB, swipe *models.Swipe) error
	FindSwipe(ctx context.Context, tx *sql.DB, swiper, swipee int) (*models.Swipe, error)
	FindAllSwipeeNotMatch(ctx context.Context, tx *sql.DB, swipee int) ([]*models.Swipe, error)
}
