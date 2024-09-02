package models

import "time"

type Message struct {
	Id       uint64
	MatchID  uint64
	SenderID uint64
	Content  string
	SendAt   time.Time
}
