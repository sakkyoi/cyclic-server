package signup

import (
	"cyclic/ent"
	"cyclic/ent/user"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (*Signup) Verify(c *gin.Context) {
	claims := c.MustGet("claims").(*magistrate.Claims) // get claims from context

	// check if the token is authorized for this action
	if !magistrate.New().Examine(claims, "verify") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "token not authorized for this action"})
		return
	}

	// get the user id from the token
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "token invalid", Detail: err.Error()})
		return
	}

	// find the user
	result, err := secretary.Minute.User.Query().
		Where(user.ID(id), user.ActiveEQ(false)). // only query inactive user, so the activated user will be regarded as not found
		Only(c)
	if ent.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ErrorResponse{Type: model.ErrorUserNotFound, Error: "user not found"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "token invalid", Detail: err.Error()})
		return
	}

	// update the user to be active
	_, err = result.Update().SetActive(true).Save(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to activate user", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: "ok"})
}
