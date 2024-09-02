package models

import (
	"encoding/json"
	"time"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type Photo struct {
	Id         uint64
	UserID     uint64
	URL        string
	IsPrimary  int8
	UploadedAt time.Time
}

type Profile struct {
	UserID           uint64
	FirstName        string
	LastName         string
	Gender           Gender
	GenderPreference time.Time
	BirthDate        time.Time
	Bio              string
	Location         string
	Interest         json.RawMessage
	Photo            []*Photo
}
