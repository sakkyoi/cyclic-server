package payment

import (
	"cyclic/pkg/magistrate"
	"cyclic/pkg/secretary"
	"cyclic/router/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type AddInput struct {
	Name    string `form:"name" binding:"required"`
	Details string `form:"details" binding:"required"` // json string
}

type Details struct {
	Method      string `json:"method" binding:"required"`
	Information string `json:"information" binding:"required"` // e.g. BankInformation
}

type BankInformation struct {
	BankName      string `json:"bank_name" binding:"required"`
	BankCode      string `json:"bank_code" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required"`
	AccountName   string `json:"account_name" binding:"required"`
}

// TODO: implement other payment methods (e.g. paypal)

func (*Payment) Add(c *gin.Context) {
	var input AddInput
	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}

	// validate json
	// parse the details
	var details Details
	if err := json.Unmarshal([]byte(input.Details), &details); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
		return
	}
	// parse the information based on the method
	switch details.Method {
	case "bank":
		var bank BankInformation
		if err := json.Unmarshal([]byte(details.Information), &bank); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: err.Error()})
			return
		}
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Type: model.ErrorInvalidInput, Error: "invalid input", Detail: "invalid method of payment"})
		return
	}

	// parse the user id from claims
	u, err := uuid.Parse(c.MustGet("claims").(*magistrate.Claims).Subject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Type: model.ErrorUnauthorized, Error: "invalid token", Detail: err.Error()})
		return
	}

	// add the payment
	p, err := secretary.Minute.Payment.Create().
		SetName(input.Name).
		SetDetails(input.Details).
		SetUserID(u).
		Save(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Type: model.ErrorInternal, Error: "internal error", Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Data: p})
}
