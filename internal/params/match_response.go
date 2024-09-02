package params

import "time"

type MatchDetailResponse struct {
	Id          uint64    `json:"id"`
	UserOne     uint64    `json:"user_one"`
	UserTwo     uint64    `json:"user_two"`
	MatchedTime time.Time `json:"matched_at"`
}
