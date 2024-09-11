package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sweatsparks/internal/models"
)

type SwipeRepository interface {
	CreateSwipe(ctx context.Context, tx *sql.Tx, swipe *models.Swipe) error
	FindSwipe(ctx context.Context, tx *sql.Tx, swiper, swipee int) (*models.Swipe, error)
	FindAllSwipeeNotMatch(ctx context.Context, tx *sql.Tx, swipee int) ([]*models.Swipe, error)
}

type SwipeRepositoryImpl struct{}

func NewSwipeRepository() SwipeRepository {
	return &SwipeRepositoryImpl{}
}

func (repository *SwipeRepositoryImpl) CreateSwipe(ctx context.Context, tx *sql.Tx, swipe *models.Swipe) error {
	sql := `INSERT INTO swipes (swiper_id, swipee_id, direction, swiped_at) values (?,?,?,?)`

	response, err := tx.ExecContext(ctx, sql, swipe.SwiperID, swipe.SwipeeID, swipe.Direction, swipe.SwipedAt)
	if err != nil {
		return err
	}
	swipeID, err := response.LastInsertId()

	if err != nil {
		return errors.New("Failed to retrieve user_id, transaction rolled back. Reason:" + err.Error())
	}

	swipe.Id = uint64(swipeID)

	return nil
}

func (repository *SwipeRepositoryImpl) FindSwipe(ctx context.Context, tx *sql.Tx, swiper, swipee int) (*models.Swipe, error) {
	sql := `SELECT swiper_id, swipee_id, direction, swiped_at FROM swipes WHERE swiper_id = ? AND swipee_id = ?`

	rows, err := tx.QueryContext(ctx, sql, swipee, swiper)
	if err != nil {
		return nil, err
	}

	var swipe models.Swipe
	if rows.Next() {
		err := rows.Scan(swipe.SwiperID, swipe.SwipeeID, swipe.Direction, swipe.SwipedAt)
		if err != nil {
			return nil, err
		}
		return &swipe, nil
	} else {
		return nil, errors.New("swipe is not found")
	}
}

func (repository *SwipeRepositoryImpl) FindAllSwipeeNotMatch(ctx context.Context, tx *sql.Tx, swipee int) ([]*models.Swipe, error) {
	sql := `SELECT id, swiper_id, swipee_id, direction, swiped_at FROM swipes WHERE swipee_id = ?`

	rows, err := tx.QueryContext(ctx, sql, swipee)
	if err != nil {
		return nil, err
	}

	var swipes []*models.Swipe
	for rows.Next() {
		var swipe models.Swipe
		if err := rows.Scan(
			swipe.Id,
			swipe.SwiperID,
			swipe.SwipeeID,
			swipe.SwipedAt,
		); err != nil {
			return nil, err
		}

		swipes = append(swipes, &swipe)
	}

	return swipes, nil
}
