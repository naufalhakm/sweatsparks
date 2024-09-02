package params

type MatchRequest struct {
	UserOne int `json:"user_one" validate:"required"`
	UserTwo int `json:"user_two" validate:"required"`
}
