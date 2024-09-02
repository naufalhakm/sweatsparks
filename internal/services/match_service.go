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

type MatchService interface {
	CreateMatchUser(ctx context.Context, req *params.MatchRequest) (*params.MatchDetailResponse, *response.CustomError)
	FindMatchDetailByUserID(ctx context.Context, userID1, UserID2 int) (*params.MatchDetailResponse, *response.CustomError)
	FindMatchAllByUserID(ctx context.Context, userID int) ([]*params.MatchDetailResponse, *response.CustomError)
}

type MatchServiceImpl struct {
	MySqlDB         *sql.DB
	MatchRepository repositories.MatchRepository
}

func NewMatchService(db *sql.DB, matchRepository repositories.MatchRepository) MatchService {
	return &MatchServiceImpl{
		MySqlDB:         db,
		MatchRepository: matchRepository,
	}
}

func (service *MatchServiceImpl) CreateMatchUser(ctx context.Context, req *params.MatchRequest) (*params.MatchDetailResponse, *response.CustomError) {
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

	_, err = service.MatchRepository.FindMatchByUserID(ctx, tx, uint64(req.UserOne), uint64(req.UserTwo))
	if err == nil {
		return nil, response.BadRequestErrorWithAdditionalInfo("User has been in match.")
	}

	var match = new(models.Match)

	match.UserOne = uint64(req.UserOne)
	match.UserTwo = uint64(req.UserTwo)
	match.MatchedTime = time.Now()

	err = service.MatchRepository.CreateMatch(ctx, tx, match)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}

	return &params.MatchDetailResponse{
		Id:          match.Id,
		UserOne:     match.UserOne,
		UserTwo:     match.UserTwo,
		MatchedTime: match.MatchedTime,
	}, nil
}

func (service *MatchServiceImpl) FindMatchDetailByUserID(ctx context.Context, userID1, UserID2 int) (*params.MatchDetailResponse, *response.CustomError) {
	tx, err := service.MySqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	result, err := service.MatchRepository.FindMatchByUserID(ctx, tx, uint64(userID1), uint64(UserID2))
	if err != nil {
		return nil, response.BadRequestErrorWithAdditionalInfo("Email has been registered.")
	}

	return &params.MatchDetailResponse{
		Id:          result.Id,
		UserOne:     result.UserOne,
		UserTwo:     result.UserTwo,
		MatchedTime: result.MatchedTime,
	}, nil
}

func (service *MatchServiceImpl) FindMatchAllByUserID(ctx context.Context, userID int) ([]*params.MatchDetailResponse, *response.CustomError) {
	tx, err := service.MySqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	matches, err := service.MatchRepository.FindAllMatchByUserID(ctx, tx, uint64(userID))
	if err != nil {
		return nil, response.BadRequestErrorWithAdditionalInfo("Email has been registered.")
	}

	var result []*params.MatchDetailResponse

	for _, match := range matches {
		result = append(result, &params.MatchDetailResponse{
			Id:          match.Id,
			UserOne:     match.UserOne,
			UserTwo:     match.UserTwo,
			MatchedTime: match.MatchedTime,
		})
	}

	return result, nil
}
