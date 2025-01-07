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
	GetUserRecomendation(ctx context.Context, userID, page, pageSize int) ([]*params.ProfileResponse, *params.Pagination, *response.CustomError)
	UpdateProfileUser(ctx context.Context, req *params.ProfileRequest) (*params.ProfileResponse, *response.CustomError)
	UpdatePhotosUser(ctx context.Context, userID int, req []*params.PhotoRequest) *response.CustomError
}

type ProfileServiceImpl struct {
	SqlDB             *sql.DB
	ProfileRepository repositories.ProfileRepository
}

func NewProfileService(db *sql.DB, profileRepository repositories.ProfileRepository) ProfileService {
	return &ProfileServiceImpl{
		SqlDB:             db,
		ProfileRepository: profileRepository,
	}
}

func (service *ProfileServiceImpl) CreateProfileUser(ctx context.Context, req *params.ProfileRequest) (*params.ProfileResponse, *response.CustomError) {
	val := validator.New()
	err := val.Struct(req)
	if err != nil {
		return nil, response.BadRequestError()
	}

	tx, err := service.SqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	birthDate, err := time.Parse("2006-01-02", req.BirthDate) // Format must match the input string format
	if err != nil {
		return nil, response.GeneralError()
	}

	preferences := models.Preference{
		PreferredAgeRange: req.Preferences.PreferredAgeRange,
		PreferredGender:   req.Preferences.PreferredGender,
		MaxDistanceKm:     req.Preferences.MaxDistanceKm,
	}

	var profile = new(models.Profile)
	profile.UserID = req.UserID
	profile.Name = req.Name
	profile.Age = req.Age
	profile.Interests = req.Interests
	profile.BirthDate = birthDate
	profile.Bio = req.Bio
	profile.Latitude = req.Latitude
	profile.Longitude = req.Longitude
	profile.LastActive = time.Now()
	profile.Preferences = &preferences

	err = service.ProfileRepository.CreateProfileByUserID(ctx, tx, profile)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}

	var photosRes []*params.PhotoResponse
	for _, ph := range req.Photos {
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
		UserID:    profile.UserID,
		Name:      profile.Name,
		Age:       profile.Age,
		Gender:    profile.Gender,
		BirthDate: profile.BirthDate,
		Bio:       profile.Bio,
		Latitude:  profile.Latitude,
		Longitude: profile.Longitude,
		Interests: profile.Interests,
		Preferences: &params.PreferenceResponse{
			MaxDistanceKm:     preferences.MaxDistanceKm,
			PreferredAgeRange: preferences.PreferredAgeRange,
			PreferredGender:   preferences.PreferredGender,
		},
		Photos: photosRes,
	}, nil
}

func (service *ProfileServiceImpl) GetProfileUser(ctx context.Context, userID int) (*params.ProfileResponse, *response.CustomError) {
	tx, err := service.SqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	result, err := service.ProfileRepository.FindProfileByUserID(ctx, tx, userID)
	if err != nil {
		return nil, response.BadRequestErrorWithAdditionalInfo("Profile not found.")
	}

	return &params.ProfileResponse{
		UserID:    result.UserID,
		Name:      result.Name,
		Age:       result.Age,
		Gender:    result.Gender,
		BirthDate: result.BirthDate,
		Bio:       result.Bio,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		Interests: result.Interests,
		Preferences: &params.PreferenceResponse{
			PreferredAgeRange: result.Preferences.PreferredAgeRange,
			PreferredGender:   result.Preferences.PreferredGender,
			MaxDistanceKm:     result.Preferences.MaxDistanceKm,
		},
	}, nil
}

func (service *ProfileServiceImpl) GetUserRecomendation(ctx context.Context, userID, page, pageSize int) ([]*params.ProfileResponse, *params.Pagination, *response.CustomError) {
	tx, err := service.SqlDB.Begin()
	if err != nil {
		return nil, nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	userLogin, err := service.ProfileRepository.FindProfileByUserID(ctx, tx, userID)
	if err != nil {
		return nil, nil, response.BadRequestErrorWithAdditionalInfo("Profile not found.")
	}

	var pagination = new(params.Pagination)
	pagination.Page = (page - 1) * pageSize
	pagination.PageSize = pageSize

	results, err := service.ProfileRepository.FindAllProfileRecomendation(ctx, tx, userLogin, pagination)
	if err != nil {
		return nil, nil, response.BadRequestErrorWithAdditionalInfo("Profile not found.")
	}

	var responses []*params.ProfileResponse
	for _, result := range results {

		photoResults, err := service.ProfileRepository.FindAllPhotoByUserID(ctx, tx, userID)
		if err != nil {
			return nil, nil, response.GeneralError()
		}

		var photoResp []*params.PhotoResponse
		for _, photoRes := range photoResults {
			photoResp = append(photoResp, &params.PhotoResponse{
				URL:       photoRes.URL,
				IsPrimary: photoRes.IsPrimary,
			})
		}

		var response = new(params.ProfileResponse)
		response.UserID = result.UserID
		response.Name = result.Name
		response.Age = result.Age
		response.Gender = result.Gender
		response.BirthDate = result.BirthDate
		response.Bio = result.Bio
		response.Latitude = result.Latitude
		response.Longitude = result.Longitude
		response.Interests = result.Interests
		response.Photos = photoResp
		response.Preferences = &params.PreferenceResponse{
			MaxDistanceKm:     result.Preferences.MaxDistanceKm,
			PreferredAgeRange: result.Preferences.PreferredAgeRange,
			PreferredGender:   result.Preferences.PreferredGender,
		}

		responses = append(responses, response)
	}

	return responses, pagination, nil
}

func (service *ProfileServiceImpl) UpdateProfileUser(ctx context.Context, req *params.ProfileRequest) (*params.ProfileResponse, *response.CustomError) {
	val := validator.New()
	err := val.Struct(req)
	if err != nil {
		return nil, response.BadRequestError()
	}

	tx, err := service.SqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	birthDate, err := time.Parse("2006-01-02", req.BirthDate) // Format must match the input string format
	if err != nil {
		return nil, response.GeneralError()
	}

	var profile = new(models.Profile)
	profile.UserID = req.UserID
	profile.Name = req.Name
	profile.Age = req.Age
	profile.Gender = req.Gender
	profile.BirthDate = birthDate
	profile.Bio = req.Bio
	profile.Latitude = req.Latitude
	profile.Longitude = req.Longitude
	profile.Interests = req.Interests

	err = service.ProfileRepository.UpdateProfileByUserID(ctx, tx, profile)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}

	return &params.ProfileResponse{
		UserID:    profile.UserID,
		Name:      profile.Name,
		Age:       profile.Age,
		Gender:    profile.Gender,
		BirthDate: profile.BirthDate,
		Bio:       profile.Bio,
		Latitude:  profile.Latitude,
		Longitude: profile.Longitude,
		Interests: profile.Interests,
		Preferences: &params.PreferenceResponse{
			PreferredAgeRange: profile.Preferences.PreferredAgeRange,
			PreferredGender:   profile.Preferences.PreferredGender,
			MaxDistanceKm:     profile.Preferences.MaxDistanceKm,
		},
	}, nil
}

func (service *ProfileServiceImpl) UpdatePhotosUser(ctx context.Context, userID int, req []*params.PhotoRequest) *response.CustomError {
	val := validator.New()
	if err := val.Struct(req); err != nil {
		return response.BadRequestError()
	}

	tx, err := service.SqlDB.Begin()
	if err != nil {
		return response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}

	defer helpers.CommitOrRollback(tx)

	for _, photo := range req {
		result, err := service.ProfileRepository.FindPhotoByID(ctx, tx, photo.Id)
		if err != nil {
			var newPhoto = new(models.Photo)
			newPhoto.UserID = uint64(userID)
			newPhoto.URL = photo.URL
			newPhoto.IsPrimary = photo.IsPrimary
			newPhoto.UploadedAt = time.Now()
			if err := service.ProfileRepository.StorePhotoByUserID(ctx, tx, newPhoto); err != nil {
				return response.GeneralError()
			}
		} else {
			var newPhoto = new(models.Photo)
			newPhoto.Id = result.Id
			newPhoto.UserID = result.UserID
			newPhoto.URL = photo.URL
			newPhoto.IsPrimary = photo.IsPrimary
			newPhoto.UploadedAt = time.Now()
			if err := service.ProfileRepository.UpdatePhotoByID(ctx, tx, newPhoto); err != nil {
				return response.GeneralError()
			}
		}
	}

	return nil
}
