package params

import "time"

type SwipeResponse struct {
	Id        uint64    `json:"id"`
	SwiperID  uint64    `json:"swiper_id"`
	SwipeeID  uint64    `json:"swipee_id"`
	Direction string    `json:"direction"`
	SwipedAt  time.Time `json:"swiped_at"`
}
