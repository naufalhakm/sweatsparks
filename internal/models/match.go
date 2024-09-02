package models

import "time"

type Match struct {
	Id          uint64
	UserOne     uint64
	UserTwo     uint64
	MatchedTime time.Time
}
