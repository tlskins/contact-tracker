package types

import "time"

// CheckInReq - request to check in user at place
type CheckInReq struct {
	PlaceID string     `json:"placeId" validate:"required"`
	In      *time.Time `json:"in" validate:"required"`
}

// CheckOutReq - request to check out user at place
type CheckOutReq struct {
	PlaceID string     `json:"placeId" validate:"required"`
	Out     *time.Time `json:"out" validate:"required"`
}
