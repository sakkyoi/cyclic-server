package signup

import (
	"cyclic/ent"
	"cyclic/pkg/colonel"
	"cyclic/pkg/dispatcher"
	"cyclic/pkg/figleaf"
	"cyclic/pkg/scribe"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Signup struct{}

type Input struct {
	Username string `form:"username" binding:"required,alphanum,min=4,max=15"`
	Password string `form:"password" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
}

func (*Signup) Signup(c *gin.Context) {
	if !colonel.Writ.Signup.Enabled {
		c.AbortWithStatusJSON(http.StatusForbidden, model.ErrorResponse{Type: model.ErrorSignupIsDisabled, Error: "signup is disabled"})
		return
	}

	var input Input
	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// hash the password
	figLeaf := figleaf.FigLeaf{}
	encoded, err := figLeaf.Cover(input.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to encode password", Detail: err.Error()})
		return
	}

	// create the user
	u := secretary.Minute.User.Create().
		SetUsername(input.Username).
		SetPassword(encoded).
		SetEmail(input.Email)

	// if verification is not required, activate the user
	if !colonel.Writ.Signup.Verification {
		u.SetActive(true)
	}

	// save the user
	result, err := u.Save(c)
	if ent.IsConstraintError(err) {
		c.AbortWithStatusJSON(http.StatusConflict, model.ErrorResponse{Type: model.ErrorUserExists, Error: "user exists"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to create user", Detail: err.Error()})
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
