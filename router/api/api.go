package api

import (
	"cyclic/router/api/plan"
	"cyclic/router/api/record"
	"cyclic/router/api/signup"
	"cyclic/router/api/subscribe"
	"cyclic/router/api/user"
)

type API struct {
	Signup    *signup.Signup
	User      *user.User
	Plan      *plan.Plan
	Subscribe *subscribe.Subscribe
	Record    *record.Record
}

func New() *API {
	return &API{}
}
