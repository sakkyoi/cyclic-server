package api

import (
	"cyclic/router/api/user"
)

type API struct {
	User *user.User
}

func New() *API {
	return &API{}
}
