package signup

import (
	"cyclic/pkg/colonel"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (*Signup) Check(c *gin.Context) {
	if !colonel.Writ.Signup.Enabled {
		c.AbortWithStatusJSON(http.StatusForbidden, model.ErrorResponse{Type: model.ErrorSignupIsDisabled, Error: "signup is disabled"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: "ok"})
}
