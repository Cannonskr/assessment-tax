package taxcalculation_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	service "github.com/Cannonskr/assessment-tax/src/api/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name           string
	requestBody    string
	expectedTax    float64
	expectedStatus int
}

func TestTaxCalculationHandler(t *testing.T) {
	e := echo.New()

	testCases := []testCase{
		{
			name:           "Case 1: Story: EXP01",
			requestBody:    `{"totalIncome": 500000, "wht": 0, "allowances": [{"allowanceType": "donation", "amount": 0}]}`,
			expectedTax:    29000.0,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Case 2: Story: EXP02",
			requestBody:    `{"totalIncome": 500000, "wht": 25000, "allowances": [{"allowanceType": "donation", "amount": 0}]}`,
			expectedTax:    4000.0,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Case 3: IStory: EXP03",
			requestBody:    `{"totalIncome": 500000, "wht": 0.0, "allowances": [{"allowanceType": "donation", "amount": 200000.0}]}`,
			expectedTax:    19000.0,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Case 4: Invalid JSON : data type",
			requestBody:    `{"totalIncome": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Case 5: No allowances",
			requestBody:    `{"totalIncome": 500000, "wht": 0, "allowances": []}`,
			expectedTax:    29000.0,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Case 6: Invalid JSON : json: unknown field \"whqt\"",
			requestBody:    `{"totalIncome": 500000, "whtt": 0, "allowances": []}`,
			expectedTax:    29000.0,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/tax/calculate", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := service.TaxCalculationHandler(c)

			assert.NoError(t, err, "No error expected")

			assert.Equal(t, tc.expectedStatus, rec.Code, "Expected status to match")

			if rec.Code == http.StatusOK {
				var result map[string]interface{}
				err := json.NewDecoder(rec.Body).Decode(&result)

				assert.NoError(t, err, "Expected JSON decoding to work")

				if expectedTax, ok := result["tax"]; ok {
					assert.Equal(t, tc.expectedTax, expectedTax, "Expected tax to match")
				}
			}
		})
	}
}
