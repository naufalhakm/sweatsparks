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

type ProfileService interface {
	CreateProfileUser(ctx context.Context, req *params.ProfileRequest) (*params.ProfileResponse, *response.CustomError)
	GetProfileUser(ctx context.Context, userID int) (*params.ProfileResponse, *response.CustomError)
	GetAllProfileUser(ctx context.Context, userID int, gender, location string) ([]*params.ProfileResponse, *response.CustomError)
	UpdateProfileUser(ctx context.Context, req *params.ProfileRequest) (*params.ProfileResponse, *response.CustomError)
}

type ProfileServiceImpl struct {
	MySqlDB           *sql.DB
	ProfileRepository repositories.ProfileRepository
}

func NewProfileService(db *sql.DB, profileRepository repositories.ProfileRepository) ProfileService {
	return &ProfileServiceImpl{
		MySqlDB:           db,
		ProfileRepository: profileRepository,
	}
}

func (service *ProfileServiceImpl) CreateProfileUser(ctx context.Context, req *params.ProfileRequest) (*params.ProfileResponse, *response.CustomError) {
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

	var profile = new(models.Profile)
	profile.UserID = req.UserID
	profile.FirstName = req.FirstName
	profile.LastName = req.LastName
	profile.Gender = req.Gender
	profile.GenderPreference = req.GenderPreference
	profile.BirthDate = req.BirthDate
	profile.Bio = req.Bio
	profile.Location = req.Location
	profile.Interest = req.Interest

	err = service.ProfileRepository.CreateProfileByUserID(ctx, tx, profile)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}

	var photosRes []*params.PhotoResponse
	for _, ph := range req.Photo {
		var photo = new(models.Photo)
		photo.UserID = req.UserID
		photo.URL = ph.URL
		photo.IsPrimary = ph.IsPrimary
		photo.UploadedAt = time.Now()

		err = service.ProfileRepository.StorePhotoByUserID(ctx, tx, photo)
		if err != nil {
			return nil, response.GeneralError(err.Error())
		}

		var photoRes = new(params.PhotoResponse)
		photoRes.URL = ph.URL
		photoRes.IsPrimary = ph.IsPrimary

		photosRes = append(photosRes, photoRes)
	}

	return &params.ProfileResponse{
		UserID:           profile.UserID,
		FirstName:        profile.FirstName,
		LastName:         profile.LastName,
		Gender:           profile.Gender,
		GenderPreference: profile.GenderPreference,
		BirthDate:        profile.BirthDate,
		Bio:              profile.Bio,
		Location:         profile.Location,
		Interest:         profile.Interest,
		Photo:            photosRes,
	}, nil
}

func (service *ProfileServiceImpl) GetProfileUser(ctx context.Context, userID int) (*params.ProfileResponse, *response.CustomError) {
	tx, err := service.MySqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	result, err := service.ProfileRepository.FindProfileByUserID(ctx, tx, userID)
	if err != nil {
		return nil, response.BadRequestErrorWithAdditionalInfo("Profile not found.")
	}

	return &params.ProfileResponse{
		UserID:           result.UserID,
		FirstName:        result.FirstName,
		LastName:         result.LastName,
		Gender:           result.Gender,
		GenderPreference: result.GenderPreference,
		BirthDate:        result.BirthDate,
		Bio:              result.Bio,
		Location:         result.Location,
		Interest:         result.Interest,
	}, nil
}

func (service *ProfileServiceImpl) GetAllProfileUser(ctx context.Context, userID int, gender, location string) ([]*params.ProfileResponse, *response.CustomError) {
	tx, err := service.MySqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	results, err := service.ProfileRepository.FindAllProfileByLocationGender(ctx, tx, location, gender)
	if err != nil {
		return nil, response.BadRequestErrorWithAdditionalInfo("Profile not found.")
	}

	var responses []*params.ProfileResponse
	for _, result := range results {
		var profile = new(params.ProfileResponse)
		profile.UserID = result.UserID
		profile.FirstName = result.FirstName
		profile.LastName = result.LastName
		profile.Gender = result.Gender
		profile.GenderPreference = result.GenderPreference
		profile.BirthDate = result.BirthDate
		profile.Bio = result.Bio
		profile.Location = result.Location
		profile.Interest = result.Interest
	}

	return responses, nil
}

func (service *ProfileServiceImpl) UpdateProfileUser(ctx context.Context, req *params.ProfileRequest) (*params.ProfileResponse, *response.CustomError) {
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

	var profile = new(models.Profile)
	profile.UserID = req.UserID
	profile.FirstName = req.FirstName
	profile.LastName = req.LastName
	profile.Gender = req.Gender
	profile.GenderPreference = req.GenderPreference
	profile.BirthDate = req.BirthDate
	profile.Bio = req.Bio
	profile.Location = req.Location
	profile.Interest = req.Interest

	err = service.ProfileRepository.UpdateProfileByUserID(ctx, tx, profile)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}

	var photosRes []*params.PhotoResponse
	for _, ph := range req.Photo {
		var photo = new(models.Photo)
		photo.UserID = req.UserID
		photo.URL = ph.URL
		photo.IsPrimary = ph.IsPrimary
		photo.UploadedAt = time.Now()

		err = service.ProfileRepository.StorePhotoByUserID(ctx, tx, photo)
		if err != nil {
			return nil, response.GeneralError(err.Error())
		}

		var photoRes = new(params.PhotoResponse)
		photoRes.URL = ph.URL
		photoRes.IsPrimary = ph.IsPrimary

		photosRes = append(photosRes, photoRes)
	}

	return &params.ProfileResponse{
		UserID:           profile.UserID,
		FirstName:        profile.FirstName,
		LastName:         profile.LastName,
		Gender:           profile.Gender,
		GenderPreference: profile.GenderPreference,
		BirthDate:        profile.BirthDate,
		Bio:              profile.Bio,
		Location:         profile.Location,
		Interest:         profile.Interest,
		Photo:            photosRes,
	}, nil
}
