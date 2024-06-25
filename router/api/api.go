package api

import (
	"cyclic/router/api/signup"
	"cyclic/router/api/user"
)

type API struct {
	Signup *signup.Signup
	User   *user.User
}

func New() *API {
	return &API{}
}
