package services

import (
	"context"
	"database/sql"
	"sweatsparks/internal/commons/response"
	"sweatsparks/internal/models"
	"sweatsparks/internal/params"
	"sweatsparks/internal/repositories"
	"sweatsparks/pkg/encryption"
	"sweatsparks/pkg/helpers"
	"sweatsparks/pkg/token"
	"time"

	"github.com/go-playground/validator"
)

type UserService interface {
	RegisterUser(ctx context.Context, req *params.UserRegisterRequest) (*params.UserRegisterResponse, *response.CustomError)
	LoginUser(ctx context.Context, req *params.UserLoginRequest) (*params.UserLoginResponse, *response.CustomError)
	GetAllUser(ctx context.Context) ([]*params.GetAllUser, *response.CustomError)
}

type UserServiceImpl struct {
	MySqlDB        *sql.DB
	UserRepository repositories.UserRepository
}

func NeewUserService(mySql *sql.DB, userRepository repositories.UserRepository) UserService {
	return &UserServiceImpl{
		MySqlDB:        mySql,
		UserRepository: userRepository,
	}
}

func (service *UserServiceImpl) RegisterUser(ctx context.Context, req *params.UserRegisterRequest) (*params.UserRegisterResponse, *response.CustomError) {
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

	_, err = service.UserRepository.FindUserByEmail(ctx, tx, req.Email)
	if err == nil {
		return nil, response.BadRequestErrorWithAdditionalInfo("Email has been registered.")
	}

	_, err = service.UserRepository.FindUserByUsername(ctx, tx, req.Username)
	if err == nil {
		return nil, response.BadRequestErrorWithAdditionalInfo("Username has been taken.")
	}

	var users = new(models.User)

	passwordHash, err := encryption.HashPassword(req.Password)
	if err != nil || passwordHash == "" {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Hashing Password Errors: %s", err.Error())
	}

	users.Email = req.Email
	users.Username = req.Username
	users.PasswordHash = passwordHash
	users.CreatedAt = time.Now()
	users.UpdatedAt = time.Now()

	err = service.UserRepository.CreateUser(ctx, tx, users)
	if err != nil {
		return nil, response.GeneralError()
	}
	token, err := token.GenerateToken(int(users.Id))
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Generate Token Errors: %s", err.Error())
	}

	response := params.UserRegisterResponse{
		Email:    users.Email,
		Username: users.Username,
		Token:    token,
	}

	return &response, nil
}

func (service *UserServiceImpl) LoginUser(ctx context.Context, req *params.UserLoginRequest) (*params.UserLoginResponse, *response.CustomError) {
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

	user, err := service.UserRepository.FindUserByEmail(ctx, tx, req.Email)
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed get user errors: %s", err.Error())
	}

	ok := encryption.CheckPasswordHash(req.Password, user.PasswordHash)
	if !ok {
		return nil, response.BadRequestErrorWithAdditionalInfo("Password wrong")
	}

	token, err := token.GenerateToken(int(user.Id))
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Generate Token Errors: %s", err.Error())
	}

	response := params.UserLoginResponse{
		Token: token,
	}

	return &response, nil
}
func (service *UserServiceImpl) GetAllUser(ctx context.Context) ([]*params.GetAllUser, *response.CustomError) {

	tx, err := service.MySqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	users, err := service.UserRepository.FindAllUser(ctx, tx)
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed get user errors: %s", err.Error())
	}

	var result []*params.GetAllUser

	for _, user := range users {
		result = append(result, &params.GetAllUser{
			ID:        user.Id,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return result, nil
}
