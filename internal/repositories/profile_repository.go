package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sweatsparks/internal/models"
	"sweatsparks/internal/params"

	"github.com/lib/pq"
)

type ProfileRepository interface {
	CreateProfileByUserID(ctx context.Context, tx *sql.Tx, profile *models.Profile) error
	FindProfileByUserID(ctx context.Context, tx *sql.Tx, userID int) (*models.Profile, error)
	FindAllProfileRecomendation(ctx context.Context, tx *sql.Tx, userLogin *models.Profile, pagination *params.Pagination) ([]*models.Profile, error)
	UpdateProfileByUserID(ctx context.Context, tx *sql.Tx, profile *models.Profile) error
	FindPhotoByID(ctx context.Context, tx *sql.Tx, id int) (*models.Photo, error)
	StorePhotoByUserID(ctx context.Context, tx *sql.Tx, photo *models.Photo) error
	FindAllPhotoByUserID(ctx context.Context, tx *sql.Tx, userID int) ([]*models.Photo, error)
	UpdatePhotoByID(ctx context.Context, tx *sql.Tx, photo *models.Photo) error
}

type ProfileRepositoryImpl struct {
}

func NewProfileRepository() ProfileRepository {
	return &ProfileRepositoryImpl{}
}

func (repository *ProfileRepositoryImpl) CreateProfileByUserID(ctx context.Context, tx *sql.Tx, profile *models.Profile) error {
	SQL := `INSERT INTO profiles (user_id,name,gender,birth_date,bio,location,interests,preferences) VALUES ($1,$2,$3,$4,$5,ST_SetSRID(ST_MakePoint($6,$7),$8,$9,$10)`

	_, err := tx.ExecContext(ctx, SQL,
		profile.UserID,
		profile.Name,
		profile.Gender,
		profile.BirthDate,
		profile.Bio,
		profile.Latitude,
		profile.Longitude,
		profile.Interests,
		profile.Preferences,
	)
	if err != nil {
		return errors.New("Failed to create a profile, transaction rolled back. Reason: " + err.Error())
	}
	return nil
}

func (repository *ProfileRepositoryImpl) FindProfileByUserID(ctx context.Context, tx *sql.Tx, userID int) (*models.Profile, error) {
	SQL := `SELECT user_id, name, age, gender, ST_X(location::geometry) AS longitude, 
           ST_Y(location::geometry) AS latitude, interests, preferences, last_active, birth_date, bio FROM profiles WHERE user_id = ?`

	rows, err := tx.QueryContext(ctx, SQL, userID)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		var profile models.Profile
		var interests []string
		var preferencesData []byte
		err := rows.Scan(&profile.UserID,
			profile.Name,
			profile.Age,
			profile.Gender,
			profile.Longitude,
			profile.Latitude,
			pq.Array(&interests),
			profile.LastActive,
			profile.BirthDate,
			profile.Bio,
		)

		var preference = new(models.Preference)
		if err := json.Unmarshal(preferencesData, preference); err != nil {
			return nil, err
		}
		profile.Preferences = preference
		if err != nil {
			return nil, err
		}
		return &profile, nil
	} else {
		return nil, errors.New("user id not found in profile")
	}
}
func (repository *ProfileRepositoryImpl) FindAllProfileRecomendation(ctx context.Context, tx *sql.Tx, userLogin *models.Profile, pagination *params.Pagination) ([]*models.Profile, error) {
	SQL := `SELECT id, name, age, gender, 
           ST_X(location::geometry) AS longitude, 
           ST_Y(location::geometry) AS latitude, 
           interests, 
           preferences,
           last_active,
           ARRAY_LENGTH(ARRAY(SELECT UNNEST(interests) INTERSECT SELECT UNNEST(ARRAY[$7::text[]])),1) AS matching_interest_count,
		   count(*) over() as total
    FROM users 
    WHERE ST_DWithin(location, ST_MakePoint($1, $2)::geography, $3 * 1000) 
    AND (age BETWEEN $4 AND $5) 
    AND gender = $6
	AND id != $8
    ORDER BY matching_interest_count ASC, last_active DESC
	LIMIT $9 OFFSET $10;`

	rows, err := tx.QueryContext(ctx, SQL,
		userLogin.Latitude,
		userLogin.Longitude,
		userLogin.Preferences.MaxDistanceKm,
		userLogin.Preferences.PreferredAgeRange[0],
		userLogin.Preferences.PreferredAgeRange[1],
		userLogin.Preferences.PreferredGender,
		userLogin.UserID,
		pagination.PageSize,
		pagination.Page,
	)
	if err != nil {
		return nil, err
	}

	var profiles []*models.Profile
	for rows.Next() {
		var interests []string
		var preferencesData []byte
		var matchingInterestCount sql.NullInt64
		var profile models.Profile
		if err := rows.Scan(&profile.UserID, &profile.Name, &profile.Age, &profile.Gender, &profile.Longitude, &profile.Latitude, pq.Array(&interests), &preferencesData, &profile.LastActive, &matchingInterestCount, &pagination.TotalCount); err != nil {
			return nil, err
		}

		var preference models.Preference
		if err := json.Unmarshal(preferencesData, &preference); err != nil {
			return nil, err
		}

		profile.Interests = interests
		profile.Preferences = &preference

		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

func (repository *ProfileRepositoryImpl) UpdateProfileByUserID(ctx context.Context, tx *sql.Tx, profile *models.Profile) error {
	SQL := `UPDATE profiles SET name = $1, age = $2, gender = $3, interests = $4, preferences = $5, birth_date = $6, bio = $6, location = ST_SetSRID(ST_MakePoint($7,$8), last_active = $9 WHERE user_id = $10`

	preferencesData, errMarshal := json.Marshal(profile.Preferences)
	if errMarshal != nil {
		return errors.New("Failed to create a profile, transaction rolled back. Reason: " + errMarshal.Error())
	}
	_, err := tx.ExecContext(ctx, SQL,
		profile.Name,
		profile.Age,
		profile.Gender,
		profile.Interests,
		preferencesData,
		profile.BirthDate,
		profile.Bio,
		profile.Latitude,
		profile.Longitude,
		profile.UserID,
	)
	if err != nil {
		return errors.New("Failed to create a profile, transaction rolled back. Reason: " + err.Error())
	}
	return nil
}

func (repository *ProfileRepositoryImpl) StorePhotoByUserID(ctx context.Context, tx *sql.Tx, photo *models.Photo) error {
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
func (repository *ProfileRepositoryImpl) FindAllPhotoByUserID(ctx context.Context, tx *sql.Tx, userID int) ([]*models.Photo, error) {
	SQL := `SELECT id, url, is_primary, uploaded_at FROM photos WHERE user_id = $1`

	rows, err := tx.QueryContext(ctx, SQL, userID)
	if err != nil {
		return nil, err
	}

	var photos []*models.Photo
	for rows.Next() {
		var photo = new(models.Photo)
		if err := rows.Scan(photo.Id, photo.URL, photo.IsPrimary, photo.UploadedAt); err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

func (repository *ProfileRepositoryImpl) UpdatePhotoByID(ctx context.Context, tx *sql.Tx, photo *models.Photo) error {
	SQL := `UPDATE photos SET url = $1, is_primary = $2, uploaded_at = $3 WHERE id = $4`

	response, err := tx.ExecContext(ctx, SQL,
		photo.URL,
		photo.IsPrimary,
		photo.UploadedAt,
		photo.Id,
	)
	if err != nil {
		return errors.New("Failed to update a photo, transaction rolled back. Reason: " + err.Error())
	}
	photoID, err := response.LastInsertId()
	if err != nil {
		return errors.New("Failed to retrieve photo_id, transaction rolled back. Reason:" + err.Error())
	}

	photo.Id = uint64(photoID)
	return nil
}

func (repository *ProfileRepositoryImpl) FindPhotoByID(ctx context.Context, tx *sql.Tx, id int) (*models.Photo, error) {
	SQL := `SELECT id, url, is_primary, uploaded_at WHERE id = $1`

	rows, err := tx.QueryContext(ctx, SQL, id)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		var photo = new(models.Photo)
		if err := rows.Scan(&photo.Id, &photo.URL, &photo.IsPrimary, &photo.UploadedAt); err != nil {
			return nil, err
		}

		return photo, nil
	} else {
		return nil, errors.New("user id not found in profile")
	}
}
