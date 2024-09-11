package params

import (
	"encoding/json"
	"time"
)

type PhotoRequest struct {
	URL       string `json:"url" validate:"required"`
	IsPrimary int8   `json:"is_primary" validate:"required"`
}

type ProfileRequest struct {
	UserID           uint64
	FirstName        string          `json:"first_name" validate:"required"`
	LastName         string          `json:"last_name" validate:"required"`
	Gender           string          `json:"gender" validate:"required"`
	GenderPreference time.Time       `json:"gender_preference" validate:"required"`
	BirthDate        time.Time       `json:"birth_date" validate:"required"`
	Bio              string          `json:"bio"`
	Location         string          `json:"location" validate:"required"`
	Interest         json.RawMessage `json:"interest" validate:"required"`
	Photo            []*PhotoRequest `json:"photo"`
}
