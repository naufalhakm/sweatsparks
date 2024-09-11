package services

import (
	"context"
	"database/sql"
	"sweatsparks/internal/commons/response"
	"sweatsparks/internal/models"
	"sweatsparks/internal/params"
	"sweatsparks/internal/repositories"
	"sweatsparks/pkg/helpers"
	"time"

	"github.com/go-playground/validator"
)

type SwipeService interface {
	CreateSwipe(ctx context.Context, req *params.SwipeRequest) (*params.SwipeResponse, *response.CustomError)
	GetSwipeBySwiperAndSwipee(ctx context.Context, swiper, swipee int) (*params.SwipeResponse, *response.CustomError)
	GetAllSwipeeNotMatchBySwipee(ctx context.Context, swipee int) ([]*params.SwipeResponse, *response.CustomError)
}

type SwipeServiceImpl struct {
	MySqlDB         *sql.DB
	SwipeRepository repositories.SwipeRepository
}

func NewSwipeRepository(db *sql.DB, swipeRepository repositories.SwipeRepository) SwipeService {
	return &SwipeServiceImpl{
		MySqlDB:         db,
		SwipeRepository: swipeRepository,
	}
}

func (service *SwipeServiceImpl) CreateSwipe(ctx context.Context, req *params.SwipeRequest) (*params.SwipeResponse, *response.CustomError) {
	val := validator.New()
	err := val.Struct(req)
	if err != nil {
		return nil, response.BadRequestError()
	}

	tx, err := service.MySqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	var swipe = new(models.Swipe)
	swipe.SwiperID = req.SwiperID
	swipe.SwipeeID = req.SwipeeID
	swipe.Direction = req.Direction
	swipe.SwipedAt = time.Now()

	err = service.SwipeRepository.CreateSwipe(ctx, tx, swipe)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}

	return &params.SwipeResponse{
		Id:        swipe.Id,
		SwiperID:  swipe.SwiperID,
		SwipeeID:  swipe.SwipeeID,
		Direction: swipe.Direction,
		SwipedAt:  swipe.SwipedAt,
	}, nil
}

func (service *SwipeServiceImpl) GetSwipeBySwiperAndSwipee(ctx context.Context, swiper, swipee int) (*params.SwipeResponse, *response.CustomError) {
	tx, err := service.MySqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	result, err := service.SwipeRepository.FindSwipe(ctx, tx, swiper, swipee)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}

	return &params.SwipeResponse{
		Id:        result.Id,
		SwiperID:  result.SwiperID,
		SwipeeID:  result.SwipeeID,
		Direction: result.Direction,
		SwipedAt:  result.SwipedAt,
	}, nil
}

func (service *SwipeServiceImpl) GetAllSwipeeNotMatchBySwipee(ctx context.Context, swipee int) ([]*params.SwipeResponse, *response.CustomError) {
	tx, err := service.MySqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	swipes, err := service.SwipeRepository.FindAllSwipeeNotMatch(ctx, tx, swipee)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}

	var result []*params.SwipeResponse
	for _, swipe := range swipes {
		result = append(result, &params.SwipeResponse{
			Id:        swipe.Id,
			SwiperID:  swipe.SwiperID,
			SwipeeID:  swipe.SwipeeID,
			Direction: swipe.Direction,
			SwipedAt:  swipe.SwipedAt,
		})
	}
	return result, nil
}
