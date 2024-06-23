package api

import (
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
)

func (a *API) Ding(c *gin.Context) {
	c.JSON(200, model.Response{Data: "dong"})
}
