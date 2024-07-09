package plan

import (
	"cyclic/ent/plan"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Input struct {
	Name         string  `form:"name" binding:"required"`
	Description  string  `form:"description" binding:"required"`
	Price        float32 `form:"price" binding:"required"`
	Currency     string  `form:"currency" binding:"required"`
	Exchangeable *bool   `form:"exchangeable" binding:"required"`
	StartFrom    string  `form:"start_from" binding:"required"`
	DurationType string  `form:"duration_type" binding:"required" enum:"days,months,years"`
	Duration     int16   `form:"duration" binding:"required"`
	AutoNotify   *bool   `form:"auto_notify" binding:"required"`
}

func (*Plan) Create(c *gin.Context) {
	var input Input
	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// parse id from claims into uuid
	hostID, err := uuid.Parse(c.MustGet("claims").(*magistrate.Claims).Subject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid token", Detail: err.Error()})
		return
	}

	// parse start_from into time
	startFrom, err := time.Parse(time.RFC3339, input.StartFrom)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "start_from must be in RFC3339 format", Detail: err.Error()})
		return
	}

	// create the plan
	result, err := secretary.Minute.Plan.Create().
		SetHostID(hostID).
		SetName(input.Name).
		SetDescription(input.Description).
		SetPrice(input.Price).
		SetCurrency(input.Currency).
		SetExchangeable(*input.Exchangeable).
		SetStartFrom(startFrom).
		SetDurationType(plan.DurationType(input.DurationType)).
		SetDuration(input.Duration).
		SetAutoNotify(*input.AutoNotify).
		Save(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to create plan", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: map[string]string{
		"id": result.ID.String(),
	}})
}
