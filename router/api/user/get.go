package user

import (
	"cyclic/pkg/magistrate"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *User) Get(c *gin.Context) {
	claims := c.MustGet("claims").(*magistrate.Claims)
	c.JSON(http.StatusOK, gin.H{"message": "get user", "claims": claims})
}
