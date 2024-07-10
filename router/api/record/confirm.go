package record

import (
	"cyclic/ent"
	"cyclic/ent/record"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (*Record) Confirm(c *gin.Context) {
	// parse the record id
	id, err := uuid.Parse(c.Param("id")) // get the record id from the path
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// get the host id from claims
	hostID, err := uuid.Parse(c.MustGet("claims").(*magistrate.Claims).Subject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "invalid token", Detail: err.Error()})
		return
	}

	// get the record
	r, err := secretary.Minute.Record.Query().
		Where(record.ID(id), record.ConfirmedAtIsNil()).
		WithSubscription().
		Only(c)
	if ent.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ErrorResponse{Type: model.ErrorRecordNotFound, Error: "record not found", Detail: err.Error()})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "internal error", Detail: err.Error()})
		return
	}

	// get plan
	plan, err := r.QuerySubscription().
		QueryPlan().
		WithHost().
		Only(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "internal error", Detail: err.Error()})
		return
	}

	// check if the host is the owner of the plan
	if plan.Edges.Host.ID != hostID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "unauthorized", Detail: "you are not the owner of the plan"})
		return
	}

	// set the confirmed_at
	result, err := r.Update().SetConfirmedAt(time.Now()).Save(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "internal error", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: result})
}
