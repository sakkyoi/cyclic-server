package record

import (
	"cyclic/ent"
	"cyclic/ent/record"
	"cyclic/ent/subscription"
	"cyclic/ent/user"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"unsafe"
)

type AddInput struct {
	Subscription string `form:"subscription" binding:"required"`
	DeclareFor   string `form:"declare_for" binding:"required"`
	Payment      string `form:"payment" binding:"required"`
	TrackingCode string `form:"tracking_code" binding:"required"`
	Remark       string `form:"remark"`
}

func (*Record) Add(c *gin.Context) {
	var input AddInput
	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// get the user id from claims
	u, err := uuid.Parse(c.MustGet("claims").(*magistrate.Claims).Subject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "invalid token", Detail: err.Error()})
		return
	}

	// parse the subscription id into uuid
	subscriptionID, err := uuid.Parse(input.Subscription)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// parse the payment id into uuid
	paymentID, err := uuid.Parse(input.Payment)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// parse the declare_for into time
	declareFor, err := time.Parse(time.RFC3339, input.DeclareFor)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// get the subscription
	s, err := secretary.Minute.Subscription.Query().
		Where(subscription.ID(subscriptionID), subscription.HasUserWith(user.ID(u)), subscription.SubscribedAtIsNil()).
		WithPlan().
		Only(c)
	if ent.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ErrorResponse{Type: model.ErrorSubscriptionNotFound, Error: "subscription not found", Detail: err.Error()})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "internal error", Detail: err.Error()})
		return
	}

	// check if the declare_for is valid
	// 1. check is the declare_for is within the subscription period
	// 2. calculate the range of the declare_for (align with the plan duration, to do poka-yoke. this must have no mistake if the request is sent by the frontend)
	// 3. check if there is no overlapping declare_for
	// check if the declare_for is within the subscription period
	if !s.StartFrom.Before(declareFor) {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: "declare_for must be after the subscription start date"})
		return
	}
	// calculate the range of the declare_for
	end := s.Edges.Plan.StartFrom
	for end.Before(declareFor) {
		// the func(b bool) *bool { return &b } is used to convert bool to *bool (get the address of the bool)
		// int(*(*byte)(unsafe.Pointer(&b))) is used to convert bool to int (the &b is made by the func above)
		end = end.AddDate(
			int(s.Edges.Plan.Duration)*int(*(*byte)(unsafe.Pointer(func(b bool) *bool { return &b }(s.Edges.Plan.DurationType == "years")))),
			int(s.Edges.Plan.Duration)*int(*(*byte)(unsafe.Pointer(func(b bool) *bool { return &b }(s.Edges.Plan.DurationType == "months")))),
			int(s.Edges.Plan.Duration)*int(*(*byte)(unsafe.Pointer(func(b bool) *bool { return &b }(s.Edges.Plan.DurationType == "days")))))
	}
	// check if there is no overlapping declare_for
	// the declare_for is the same as the end
	records, err := secretary.Minute.Record.Query().
		Where(record.HasSubscriptionWith(subscription.ID(subscriptionID)), record.DeclareForEQ(end)).
		All(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "internal error", Detail: err.Error()})
		return
	}
	if len(records) != 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: "declare_for is overlapping with existing record"})
		return
	}

	// create the record
	result, err := secretary.Minute.Record.Create().
		SetSubscription(s).
		SetDeclareFor(end).
		SetPaymentID(paymentID).
		SetTrackingCode(input.TrackingCode).
		SetRemark(input.Remark).
		Save(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "internal error", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: result})
}
