package types

// SignInReq - request to check in place at place
type SignInReq struct {
	Password string `json:"password" validate:"required"`
}
