package subscribe

import (
	"cyclic/ent"
	"cyclic/ent/subscription"
	"cyclic/ent/user"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Subscribe struct{}

func (*Subscribe) Subscribe(c *gin.Context) {
	// parse the subscription id
	id, err := uuid.Parse(c.Param("id")) // get the plan id from the path
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// get the user id from claims
	u, err := uuid.Parse(c.MustGet("claims").(*magistrate.Claims).Subject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "invalid token", Detail: err.Error()})
		return
	}

	// get the subscription
	s, err := secretary.Minute.Subscription.Query().
		Where(subscription.ID(id), subscription.HasUserWith(user.ID(u)), subscription.SubscribedAtIsNil()).
		Only(c)
	if ent.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ErrorResponse{Type: model.ErrorSubscriptionNotFound, Error: "subscription not found", Detail: err.Error()})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "internal error", Detail: err.Error()})
		return
	}

	// set the subscribed_at
	result, err := s.Update().SetSubscribedAt(time.Now()).Save(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "internal error", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: result})
}
