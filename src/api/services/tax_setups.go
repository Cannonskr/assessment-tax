package taxcalculation

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AllowanceInput struct {
	Amount float64 `json:"amount"`
}

type PersonalDeductionResponse struct {
	Personal float64 `json:"personalDeduction"`
}

type KReceiptDeductionResponse struct {
	KReceipt float64 `json:"kReceipt"`
}

func UpdatePersonalDeductionHandler(c echo.Context) error {
	var input AllowanceInput

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if input.Amount < 10000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Amount cannot be less than 10,000"})
	}

	if input.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Amount cannot be more than 100,000"})
	}

	allowedAllowanceTypes["personal"] = input.Amount

	response := PersonalDeductionResponse{
		Personal: input.Amount,
	}

	return c.JSON(http.StatusOK, response)
}

func UpdateKReceiptDeductionHandler(c echo.Context) error {
	var input AllowanceInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if input.Amount < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Amount cannot be less than 0"})
	}

	if input.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Amount cannot be more than 100,000"})
	}

	allowedAllowanceTypes["k-receipt"] = input.Amount

	response := KReceiptDeductionResponse{
		KReceipt: input.Amount,
	}

	return c.JSON(http.StatusOK, response)
}
