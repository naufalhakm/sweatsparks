package params

type SwipeRequest struct {
	SwiperID  uint64 `json:"swiper_id"`
	SwipeeID  uint64 `json:"swipee_id"`
	Direction string `json:"direction"`
}
