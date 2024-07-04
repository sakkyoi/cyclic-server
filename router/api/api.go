package api

import (
	"cyclic/router/api/plan"
	"cyclic/router/api/signup"
	"cyclic/router/api/user"
)

type API struct {
	Signup *signup.Signup
	User   *user.User
	Plan   *plan.Plan
}

func New() *API {
	return &API{}
}
