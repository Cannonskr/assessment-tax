package taxcalculation

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TaxData struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

var allowedAllowanceTypes = map[string]float64{
	"donation":  100000,
	"k-receipt": 50000,
	"personal":  60000,
}

type TaxBracket struct {
	MinIncome   float64
	MaxIncome   float64
	Rate        float64
	Description string
}

type TaxResult struct {
	TotalTax  float64    `json:"tax"`
	TaxRefund float64    `json:"taxRefund,omitempty"`
	TaxLevel  []TaxLevel `json:"taxLevel,omitempty"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

func validateAndAdjustAllowance(allowance *Allowance) error {
	limit, exists := allowedAllowanceTypes[allowance.AllowanceType]
	if !exists {
		return errors.New("invalid allowance type")
	}

	if allowance.Amount > limit {
		allowance.Amount = limit
	}

	return nil
}

func checkDuplicateAllowanceTypes(allowances []Allowance) error {
	allowanceTypes := make(map[string]bool)

	for _, allowance := range allowances {

		if allowance.AllowanceType == "personal" {
			return errors.New("not allowed type:" + allowance.AllowanceType)
		}
		if _, exists := allowanceTypes[allowance.AllowanceType]; exists {
			return errors.New("duplicate allowance type: " + allowance.AllowanceType)
		}
		allowanceTypes[allowance.AllowanceType] = true
	}

	return nil
}

func calculateTax(data TaxData) TaxResult {

	taxableIncome := data.TotalIncome

	for _, allowance := range data.Allowances {
		taxableIncome -= allowance.Amount
	}

	taxBrackets := []TaxBracket{
		{MinIncome: 0, MaxIncome: 150000, Rate: 0.00, Description: "0-150,000"},
		{MinIncome: 150001, MaxIncome: 500000, Rate: 0.10, Description: "150,001-500,000"},
		{MinIncome: 500001, MaxIncome: 1000000, Rate: 0.15, Description: "500,001-1,000,000"},
		{MinIncome: 1000001, MaxIncome: 2000000, Rate: 0.20, Description: "1,000,001-2,000,000"},
		{MinIncome: 2000001, MaxIncome: math.MaxFloat64, Rate: 0.35, Description: "2,000,001 ขึ้นไป"},
	}

	totalTax := 0.0
	taxLevels := []TaxLevel{}

	for _, bracket := range taxBrackets {

		if taxableIncome > bracket.MinIncome {
			maxTaxable := taxableIncome
			if maxTaxable > bracket.MaxIncome {
				maxTaxable = bracket.MaxIncome
			}
			taxable := (maxTaxable - (bracket.MinIncome - 1))
			tax := taxable * bracket.Rate

			if tax >= data.WHT {
				tax -= data.WHT
				data.WHT = 0
			} else {
				data.WHT -= tax
				tax = 0
			}

			totalTax += tax

			taxLevels = append(taxLevels, TaxLevel{
				Level: bracket.Description,
				Tax:   math.Round(tax*10) / 10,
			})
		} else {
			taxLevels = append(taxLevels, TaxLevel{
				Level: bracket.Description,
				Tax:   0,
			})
		}
	}

	totalTax -= data.WHT
	taxRefund := 0.0
	if totalTax < 0 {
		taxRefund = -totalTax
		totalTax = 0.0
	}

	totalTax = math.Round(totalTax*10) / 10
	taxRefund = math.Round(taxRefund*10) / 10

	return TaxResult{
		TotalTax:  totalTax,
		TaxRefund: taxRefund,
		TaxLevel:  taxLevels,
	}
}

func TaxCalculationHandler(c echo.Context) error {

	var data TaxData
	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON: " + err.Error()})
	}

	if err := checkDuplicateAllowanceTypes(data.Allowances); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	for i := range data.Allowances {
		if err := validateAndAdjustAllowance(&data.Allowances[i]); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}

	amount := allowedAllowanceTypes["personal"]
	personalAllowance := Allowance{
		AllowanceType: "personal",
		Amount:        amount,
	}

	data.Allowances = append(data.Allowances, personalAllowance)

	result := calculateTax(data)
	return c.JSON(http.StatusOK, result)

}
