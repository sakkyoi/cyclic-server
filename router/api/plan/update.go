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

type UpdateInput struct {
	Name         string  `form:"name"`
	Description  string  `form:"description"`
	Price        float32 `form:"price"`
	Currency     string  `form:"currency"`
	Exchangeable *bool   `form:"exchangeable"`
	StartFrom    string  `form:"start_from"`
	DurationType string  `form:"duration_type" enum:"days,months,years"`
	Duration     int16   `form:"duration"`
	AutoNotify   *bool   `form:"auto_notify"`
}

func (*Plan) Update(c *gin.Context) {
	var input UpdateInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// parse plan id from params into uuid
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// parse id from claims into uuid
	hostID, err := uuid.Parse(c.MustGet("claims").(*magistrate.Claims).Subject)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "invalid token", Detail: err.Error()})
		return
	}

	// update the plan
	query := secretary.Minute.Plan.UpdateOneID(planID).
		Where(plan.DeletedAtIsNil(), plan.HasHostWith(user.ID(hostID)))

	// set fields if they are provided in the input
	// 0 means nil for float32 and int16
	if input.Name != "" {
		query.SetName(input.Name)
	}

	if input.Description != "" {
		query.SetDescription(input.Description)
	}

	if input.Price != 0 {
		query.SetPrice(input.Price)
	}

	if input.Currency != "" {
		query.SetCurrency(input.Currency)
	}

	if input.Exchangeable != nil {
		query.SetExchangeable(*input.Exchangeable)
	}

	if input.StartFrom != "" {
		startFrom, err := time.Parse(time.RFC3339, input.StartFrom) // parse start_from into time
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
			return
		}
		query.SetStartFrom(startFrom)
	}

	if input.DurationType != "" {
		query.SetDurationType(plan.DurationType(input.DurationType))
	}

	if input.Duration != 0 {
		query.SetDuration(input.Duration)
	}

	if input.AutoNotify != nil {
		query.SetAutoNotify(*input.AutoNotify)
	}

	// save the plan
	result, err := query.Save(c)
	if ent.IsNotFound(err) {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Type: model.ErrorPlanNotFound, Error: "plan not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to update plan", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: map[string]string{
		"id": result.ID.String(),
	}})
}
