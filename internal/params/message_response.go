package params

import "time"

type MessageResponse struct {
	Id       uint64    `json:"id"`
	MatchID  uint64    `json:"match_id"`
	SenderID uint64    `json:"sender_id"`
	Content  string    `json:"content"`
	SendAt   time.Time `json:"sent_at"`
}
