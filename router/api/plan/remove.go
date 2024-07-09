package plan

import (
	"cyclic/ent"
	"cyclic/ent/plan"
	"cyclic/ent/user"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (*Plan) Remove(c *gin.Context) {
	// parse the plan id
	id, err := uuid.Parse(c.Param("id")) // get the plan id from the path
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// parse id from claims into uuid
	hostID, err := uuid.Parse(c.MustGet("claims").(*magistrate.Claims).Subject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "invalid token", Detail: err.Error()})
		return
	}

	// soft delete the plan (by setting DeletedAt)
	_, err = secretary.Minute.Plan.Update().
		Where(plan.IDEQ(id), plan.DeletedAtIsNil(), plan.HasHostWith(user.ID(hostID))).
		SetDeletedAt(time.Now()).
		Save(c)
	if ent.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ErrorResponse{Type: model.ErrorPlanNotFound, Error: "plan not found"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to remove plan", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: "ok"})
}
