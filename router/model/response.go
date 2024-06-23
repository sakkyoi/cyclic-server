package model

var (
	ErrorSignupIsDisabled = "signup_is_disabled"
	ErrorInvalidInput     = "invalid_input"
	ErrorInternal         = "internal"
	ErrorUserExists       = "user_exists"
	ErrorUnauthorized     = "unauthorized"
	ErrorUserNotFound     = "user_not_found"
	ErrorUserNotActivated = "user_not_activated"
)

type Response struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Type   string `json:"type"`
	Error  string `json:"error"`
	Detail string `json:"detail"`
}
