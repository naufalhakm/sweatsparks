package params

import (
	"encoding/json"
	"time"
)

type PhotoResponse struct {
	URL       string `json:"url"`
	IsPrimary int8   `json:"is_primary"`
}

type ProfileResponse struct {
	UserID           uint64           `json:"user_id"`
	FirstName        string           `json:"first_name"`
	LastName         string           `json:"last_name"`
	Gender           string           `json:"gender"`
	GenderPreference time.Time        `json:"gender_preference"`
	BirthDate        time.Time        `json:"birth_date"`
	Bio              string           `json:"bio"`
	Location         string           `json:"location"`
	Interest         json.RawMessage  `json:"interest"`
	Photo            []*PhotoResponse `json:"photo"`
}
