package model

var (
	ErrorSignupIsDisabled     = "signup_is_disabled"
	ErrorInvalidInput         = "invalid_input"
	ErrorInternal             = "internal"
	ErrorUserExists           = "user_exists"
	ErrorUnauthorized         = "unauthorized"
	ErrorUserNotFound         = "user_not_found"
	ErrorUserNotActivated     = "user_not_activated"
	ErrorPlanNotFound         = "plan_not_found"
	ErrorSubscriptionExists   = "subscription_exists"
	ErrorSubscriptionNotFound = "subscription_not_found"
)

type Response struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Type   string `json:"type"`
	Error  string `json:"error"`
	Detail string `json:"detail"`
}
