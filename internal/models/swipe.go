package models

import "time"

type Swipe struct {
	Id        uint64
	SwiperID  uint64
	SwipeeID  uint64
	Direction string
	SwipedAt  time.Time
}
