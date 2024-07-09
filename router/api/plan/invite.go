package plan

import (
	"cyclic/ent"
	"cyclic/ent/plan"
	"cyclic/ent/user"
	"cyclic/pkg/dispatcher"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/scribe"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type InviteInput struct {
	User string `form:"user" binding:"required"`
}

func (*Plan) Invite(c *gin.Context) {
	var input InviteInput
	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

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

	// parse the user id into uuid
	userID, err := uuid.Parse(input.User)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid user id", Detail: err.Error()})
		return
	}

	// check if the user is the host
	if hostID == userID {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "cannot invite yourself"})
		return
	}

	// get the plan (with the host)
	p, err := secretary.Minute.Plan.Query().
		Where(plan.ID(id), plan.HasHostWith(user.ID(hostID))).
		Only(c)
	if ent.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ErrorResponse{Type: model.ErrorPlanNotFound, Error: "plan not found"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to get plan", Detail: err.Error()})
		return
	}

	// get the user
	u, err := secretary.Minute.User.Query().
		Where(user.ID(userID)).
		Only(c)
	if ent.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ErrorResponse{Type: model.ErrorUserNotFound, Error: "user not found"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to get user", Detail: err.Error()})
		return
	}

	// invite the user to the plan
	// the reason we don't use AddUserID is that we need to check if the user is existed.
	// user already invited error are the same as user exists error (constraint error)
	result, err := secretary.Minute.Subscribe.Create().
		SetPlan(p).
		SetUser(u).
		Save(c)
	if ent.IsConstraintError(err) {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorSubscriptionExists, Error: "user already invited to the plan"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "failed to invite user to the plan", Detail: err.Error()})
		return
	}

	// enqueue a message to send an email
	if err := dispatcher.Enqueue(&dispatcher.Message{
		Type:   dispatcher.Invite,
		Target: u.ID.String(),
		Data:   p.ID.String(),
	}); err != nil {
		scribe.Scribe.Error("failed to enqueue message", zap.Error(err)) // just log the error cause the user is already created
	}

	c.JSON(http.StatusOK, model.Response{Data: result.ID.String()})
}
