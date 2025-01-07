package params

import (
	"time"
)

type PhotoResponse struct {
	URL       string `json:"url"`
	IsPrimary int8   `json:"is_primary"`
}

type PreferenceResponse struct {
	PreferredAgeRange [2]int `json:"age"`
	PreferredGender   string `json:"gender"`
	MaxDistanceKm     int    `json:"distance"`
}

type ProfileResponse struct {
	UserID      uint64              `json:"user_id"`
	Name        string              `json:"name"`
	Age         int32               `json:"age"`
	Gender      string              `json:"gender"`
	BirthDate   time.Time           `json:"birth_date"`
	Bio         string              `json:"bio"`
	Latitude    float64             `json:"latitude"`
	Longitude   float64             `json:"longitude"`
	Interests   []string            `json:"interests"`
	Photos      []*PhotoResponse    `json:"photos"`
	Preferences *PreferenceResponse `json:"preferences"`
}
