package user

import (
	"cyclic/ent"
	"cyclic/ent/user"
	"cyclic/pkg/figleaf"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthInput struct {
	Username string `form:"username" binding:"required,alphanum,min=4,max=15"`
	Password string `form:"password" binding:"required"`
}

func (*User) Auth(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	figLeaf := figleaf.FigLeaf{}

	// query user
	result, err := secretary.Minute.User.Query().
		Where(user.Username(input.Username)).
		Only(c)
	if ent.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUserNotFound, Error: "user not found"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to query user", Detail: err.Error()})
		return
	}

	// verify is password correct
	ok, err := figLeaf.Peep(input.Password, result.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to encode password", Detail: err.Error()})
		return
	} else if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "unauthorized"})
		return
	}

	// check if user is activated
	if !result.Active {
		c.AbortWithStatusJSON(http.StatusForbidden, model.ErrorResponse{Type: model.ErrorUserNotActivated, Error: "user not activated"})
		return
	}

	// issue token
	m := magistrate.New()
	token, err := m.Issue([]string{"general"}, result.ID.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to issue token", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: map[string]interface{}{"token": token}})
}
