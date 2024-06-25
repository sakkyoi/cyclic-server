package signup

import (
	"cyclic/ent"
	"cyclic/ent/user"
	"cyclic/pkg/dispatcher"
	"cyclic/pkg/scribe"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ResendInput struct {
	Email string `form:"email" binding:"required,email"`
}

func (s *Signup) Resend(c *gin.Context) {
	var input ResendInput
	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// query user by email
	result, err := secretary.Minute.User.Query().
		Where(user.Email(input.Email), user.ActiveEQ(false)). // only query inactive user, so the activated user will be regarded as not found
		Only(c)
	if ent.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ErrorResponse{Type: model.ErrorUserNotFound, Error: "user not found"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to query user", Detail: err.Error()})
		return
	}

	// enqueue a message to send an email
	if err := dispatcher.Enqueue(&dispatcher.Message{
		Type:   dispatcher.Verify,
		Target: result.ID.String(),
	}); err != nil {
		scribe.Scribe.Error("failed to enqueue message", zap.Error(err)) // just log the error cause the user is already created
	}

	c.JSON(http.StatusOK, model.Response{Data: "ok"})
}
