package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
	"github.com/paulantezana/requirement/utilities"
	"net/http"
	"time"
)

type quotationDetailResponse struct {
	ID             uint    `json:"id"`
	Amount         float32 `json:"amount"`
	UnitMeasure    string  `json:"unit_measure"`
	ProductID      uint    `json:"product_id"`
	ProductName    string  `json:"product_name"`
	SuggestedPrice float32 `json:"suggested_price"`
	UnitPrice      float32 `json:"unit_price"`
	Observation    string  `json:"observation"`
}

func GetRequireByRequirement(c echo.Context) error {
	// Get data request
	require := models.Require{}
	if err := c.Bind(&require); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Find in database requires
	quotationDetailResponses := make([]quotationDetailResponse, 0)
	if err := db.Table("requires").
		Select("requires.id, requires.amount, requires.unit_measure, products.name as product_name, requires.product_id, requires.suggested_price, requires.observation").
		Joins("INNER JOIN products on requires.product_id = products.id").
		Where("requires.requirement_id = ?", require.RequirementID).
		Scan(&quotationDetailResponses).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quotationDetailResponses,
	})
}

type quotationResponse struct {
	ID            uint      `json:"id"`
	EmissionDate  time.Time `json:"emission_date"`
	SuggestWinner bool      `json:"suggest_winner"` // Winner suggestion by user
	Observation   string    `json:"observation"`

	ProviderID    uint   `json:"provider_id"`
	ProviderName  string `json:"provider_name"`
	RequirementID uint   `json:"requirement_id"`

	QuotationDetails []quotationDetailResponse `json:"quotation_details"`
}

func GetRequireByQuotation(c echo.Context) error {
	// Get data request
	quotation := models.Quotation{}
	if err := c.Bind(&quotation); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Find in database requires
	quotationDetailResponses := make([]quotationDetailResponse, 0)
	if err := db.Table("quotations").
		Select("quotation_details.id, quotation_details.quotation_id, requires.amount, requires.unit_measure, products.name, requires.suggested_price, quotation_details.unit_price, requires.observation").
		Joins("INNER JOIN quotation_details on quotations.id = quotation_details.quotation_id").
		Joins("INNER JOIN requires on quotation_details.require_id = requires.id").
		Joins("INNER JOIN products on requires.product_id = products.id").
		Where("quotations.id = ?", quotation.ID).
		Scan(&quotationDetailResponses).Error; err != nil {
		return err
	}

	// Find quotation
	quotationRes := make([]quotationResponse, 0)
	if err := db.Table("quotations").
		Select("quotations.id, quotations.emission_date, quotations.suggest_winner, quotations.observation, quotations.provider_id, providers.name as provider_name, quotations.requirement_id").
		Joins("INNER JOIN providers on quotations.provider_id = providers.id").
		Where("quotations.id = ?", quotation.ID).
		Scan(&quotationRes).Error; err != nil {
		return err
	}

	// Customise response
	quotationData := quotationRes[0]
	quotationData.QuotationDetails = quotationDetailResponses

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    quotationData,
	})
}

func DeleteRequire(c echo.Context) error {
	// Get data request
	require := models.Require{}
	if err := c.Bind(&require); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation product exist
	if db.First(&require).RecordNotFound() {
		return c.JSON(http.StatusCreated, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", require.ID),
		})
	}

	// Delete product in database
	if err := db.Delete(&require).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    require.ID,
	})
}
