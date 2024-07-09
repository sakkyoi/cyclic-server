package subscribe

import (
	"cyclic/ent/subscribe"
	"cyclic/ent/user"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (*Subscribe) List(c *gin.Context) {
	// parse id from claims into uuid
	id, err := uuid.Parse(c.MustGet("claims").(*magistrate.Claims).Subject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "invalid token", Detail: err.Error()})
		return
	}

	// get subscriptions
	subscriptions, err := secretary.Minute.Subscribe.Query().
		Where(subscribe.HasUserWith(user.ID(id))).
		WithPlan().
		All(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "internal error", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: subscriptions})
}
