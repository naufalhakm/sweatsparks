package models

import "time"

type Swipe struct {
	SwiperID  uint64
	SwipeeID  uint64
	Direction string
	SwipedAt  time.Time
}
