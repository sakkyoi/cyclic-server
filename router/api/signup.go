package api

import (
	"cyclic/ent"
	"cyclic/ent/user"
	"cyclic/pkg/colonel"
	"cyclic/pkg/dispatcher"
	"cyclic/pkg/figleaf"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/scribe"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type SignupInput struct {
	Username string `form:"username" binding:"required,alphanum,min=4,max=32"`
	Password string `form:"password" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
}

func (a *API) CheckIsSignupEnabled(c *gin.Context) {
	if !colonel.Writ.Signup.Enabled {
		c.AbortWithStatusJSON(http.StatusForbidden, model.ErrorResponse{Type: model.ErrorSignupIsDisabled, Error: "signup is disabled"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: "ok"})
}

func (a *API) Signup(c *gin.Context) {
	if !colonel.Writ.Signup.Enabled {
		c.AbortWithStatusJSON(http.StatusForbidden, model.ErrorResponse{Type: model.ErrorSignupIsDisabled, Error: "signup is disabled"})
		return
	}

	var input SignupInput
	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	figLeaf := figleaf.FigLeaf{}
	encoded, err := figLeaf.Cover(input.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to encode password", Detail: err.Error()})
		return
	}

	u := secretary.Minute.User.Create().
		SetUsername(input.Username).
		SetPassword(encoded).
		SetEmail(input.Email)

	if !colonel.Writ.Signup.Verification {
		u.SetActive(true)
	}

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
		Type:   dispatcher.Signup,
		Target: result.ID.String(),
	}); err != nil {
		scribe.Scribe.Error("failed to enqueue message", zap.Error(err)) // just log the error cause the user is already created
	}

	c.JSON(http.StatusOK, model.Response{Data: "ok"})
}

func (a *API) VerifySignup(c *gin.Context) {
	claims := c.MustGet("claims").(*magistrate.Claims) // get claims from context

	// check if the token is authorized for this action
	if !magistrate.New().Examine(claims, "signup") {
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
		Where(user.ID(id), user.ActiveEQ(false)).
		Only(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "token invalid", Detail: err.Error()})
		return
	}

	_, err = result.Update().SetActive(true).Save(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to activate user", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: "ok"})
}
