package payment

import (
	"cyclic/ent/payment"
	"cyclic/ent/user"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (*Payment) Remove(c *gin.Context) {
	// parse the payment id
	id, err := uuid.Parse(c.Param("id")) // get the plan id from the path
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// parse id from claims into uuid
	userID, err := uuid.Parse(c.MustGet("claims").(*magistrate.Claims).Subject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "invalid token", Detail: err.Error()})
		return
	}

	// remove the payment
	_, err = secretary.Minute.Payment.Delete().
		Where(payment.ID(id), payment.HasUserWith(user.ID(userID))).
		Exec(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to remove payment", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: "ok"})
}
