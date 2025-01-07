package models

import (
	"time"
)

type Photo struct {
	Id         uint64
	UserID     uint64
	URL        string
	IsPrimary  int8
	UploadedAt time.Time
}

type Preference struct {
	PreferredAgeRange [2]int `json:"age"`
	PreferredGender   string `json:"gender"`
	MaxDistanceKm     int    `json:"distance"`
}

type Profile struct {
	UserID      uint64
	Name        string
	Age         int32
	Gender      string
	Interests   []string
	Preferences *Preference
	BirthDate   time.Time
	Bio         string
	Latitude    float64
	Longitude   float64
	Photos      []*Photo
	LastActive  time.Time
}
