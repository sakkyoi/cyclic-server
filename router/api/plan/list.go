package plan

import (
	"cyclic/ent/plan"
	"cyclic/ent/user"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (*Plan) List(c *gin.Context) {
	// parse id from claims into uuid
	userId, err := uuid.Parse(c.MustGet("claims").(*magistrate.Claims).Subject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "invalid token", Detail: err.Error()})
		return
	}

	// list plans
	result, err := secretary.Minute.Plan.Query().
		Where(plan.HasHostWith(user.ID(userId)), plan.DeletedAtIsNil()).
		All(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to list plans", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: result})
}
