package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sweatsparks/internal/models"
)

type ProfileRepository interface {
	CreateProfileByUserID(ctx context.Context, tx *sql.DB, profile *models.Profile) error
	FindProfileByUserID(ctx context.Context, tx *sql.DB, userID int) (*models.Profile, error)
	FindAllProfileByLocationGender(ctx context.Context, tx *sql.DB, location, gender string) ([]*models.Profile, error)
	UpdateProfileByUserID(ctx context.Context, tx *sql.DB, profile *models.Profile) error
	StorePhotoByUserID(ctx context.Context, tx *sql.DB, photo *models.Photo) error
}

type ProfileRepositoryImpl struct {
}

func NewProfileRepository() ProfileRepository {
	return &ProfileRepositoryImpl{}
}

func (repository *ProfileRepositoryImpl) CreateProfileByUserID(ctx context.Context, tx *sql.DB, profile *models.Profile) error {
	SQL := `INSERT INTO profiles (user_id,first_name,last_name,gender,gender_preference,date_of_birth,bio,location,interests) VALUES (?,?,?,?,?,?,?,?,?,?,?)`

	_, err := tx.ExecContext(ctx, SQL,
		profile.UserID,
		profile.FirstName,
		profile.LastName,
		profile.Gender,
		profile.GenderPreference,
		profile.BirthDate,
		profile.Bio,
		profile.Location,
		profile.Interest,
	)
	if err != nil {
		return errors.New("Failed to create a profile, transaction rolled back. Reason: " + err.Error())
	}
	return nil
}

func (repository *ProfileRepositoryImpl) FindProfileByUserID(ctx context.Context, tx *sql.DB, userID int) (*models.Profile, error) {
	SQL := `SELECT user_id,first_name,last_name,gender,gender_preference,date_of_birth,bio,location,interests FROM profiles WHERE user_id = ?`

	rows, err := tx.QueryContext(ctx, SQL, userID)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		var profile models.Profile
		err := rows.Scan(&profile.UserID,
			profile.FirstName,
			profile.LastName,
			profile.Gender,
			profile.GenderPreference,
			profile.BirthDate,
			profile.Bio,
			profile.Location,
			profile.Interest,
		)
		if err != nil {
			return nil, err
		}
		return &profile, nil
	} else {
		return nil, errors.New("user id not found in profile")
	}
}
func (repository *ProfileRepositoryImpl) FindAllProfileByLocationGender(ctx context.Context, tx *sql.DB, location, gender string) ([]*models.Profile, error) {
	SQL := `SELECT user_id,first_name,last_name,gender,gender_preference,date_of_birth,bio,location,interests FROM profiles WHERE location = ? AND gender = ?`

	rows, err := tx.QueryContext(ctx, SQL, location, gender)
	if err != nil {
		return nil, err
	}

	var profiles []*models.Profile
	for rows.Next() {
		var profile models.Profile
		if err := rows.Scan(&profile.UserID,
			profile.FirstName,
			profile.LastName,
			profile.Gender,
			profile.GenderPreference,
			profile.BirthDate,
			profile.Bio,
			profile.Location,
			profile.Interest,
		); err != nil {
			return nil, err
		}
		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

func (repository *ProfileRepositoryImpl) UpdateProfileByUserID(ctx context.Context, tx *sql.DB, profile *models.Profile) error {
	SQL := `UPDATE profiles SET first_name = ?,last_name = ?,gender = ?,gender_preference = ?,date_of_birth = ?,bio =?,location = ?,interests = ? WHERE user_id = ?`

	_, err := tx.ExecContext(ctx, SQL,
		profile.FirstName,
		profile.LastName,
		profile.Gender,
		profile.GenderPreference,
		profile.BirthDate,
		profile.Bio,
		profile.Location,
		profile.Interest,
		profile.UserID,
	)
	if err != nil {
		return errors.New("Failed to create a profile, transaction rolled back. Reason: " + err.Error())
	}
	return nil
}

func (repository *ProfileRepositoryImpl) StorePhotoByUserID(ctx context.Context, tx *sql.DB, photo *models.Photo) error {
	SQL := `INSERT INTO photos (user_id,url,is_primary,uploaded_at) VALUES (?,?,?,?)`

	response, err := tx.ExecContext(ctx, SQL,
		photo.UserID,
		photo.URL,
		photo.IsPrimary,
		photo.UploadedAt,
	)
	if err != nil {
		return errors.New("Failed to create a photo, transaction rolled back. Reason: " + err.Error())
	}
	photoID, err := response.LastInsertId()
	if err != nil {
		return errors.New("Failed to retrieve photo_id, transaction rolled back. Reason:" + err.Error())
	}

	photo.Id = uint64(photoID)
	return nil
}
