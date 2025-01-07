package params

type PhotoRequest struct {
	Id        int    `json:"id"`
	URL       string `json:"url" validate:"required"`
	IsPrimary int8   `json:"is_primary" validate:"required"`
}

type PreferenceRequest struct {
	PreferredAgeRange [2]int `json:"age"`
	PreferredGender   string `json:"gender"`
	MaxDistanceKm     int    `json:"distance"`
}

type ProfileRequest struct {
	UserID      uint64
	Name        string             `json:"name" validate:"required"`
	Age         int32              `json:"age" validate:"required"`
	Gender      string             `json:"gender" validate:"required"`
	BirthDate   string             `json:"birth_date" validate:"required"`
	Bio         string             `json:"bio"`
	Latitude    float64            `json:"latitude" validate:"required"`
	Longitude   float64            `json:"longitude" validate:"required"`
	Interests   []string           `json:"interests" validate:"required"`
	Photos      []PhotoRequest     `json:"photos"`
	Preferences *PreferenceRequest `json:"preferences" validate:"required"`
}
