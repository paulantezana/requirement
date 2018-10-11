package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
	"github.com/paulantezana/requirement/utilities"
	"net/http"
	"time"
)

type quotationCustomResponse struct {
	ID            uint    `json:"id"`
	ProviderID    uint    `json:"provider_id"`
	ProviderName  string  `json:"provider_name"`
	UserID        uint    `json:"user_id"`
	UserFirstName string  `json:"user_first_name"`
	UserLastName  string  `json:"user_last_name"`
	RequirementID uint    `json:"requirement_id"`
	Count         uint    `json:"count"`
	WinnerLevel   uint    `json:"winner_level"`
	Winner        bool    `json:"winner"`
	Summation     float32 `json:"summation"`
}

type quotationResult struct {
	ID            uint
	ProviderID    uint
	ProviderName  string
	UserID        uint
	UserFirstName string
	UserLastName  string
	RequirementID uint
	Count         uint
	WinnerLevel   uint
	Winner        bool
}

type winnerLevelResult struct {
	ID            uint
	ProviderID    uint
	RequirementID uint
	Summation     float32
	SuggestWinner bool
}

func GetQuotations(c echo.Context) error {
	// Get data request
	request := utilities.RequestQuotation{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Find quotations in database by RequirementID  ========== Quotations, Providers, Users
	quotationResults := make([]quotationResult, 0)
	if err := db.Table("quotations").
		Select("quotations.id, providers.id as provider_id, providers.name as provider_name, users.id as user_id, users.first_name as user_first_name, users.last_name as user_last_name, quotations.requirement_id, count(*), quotations.winner_level, quotations.Winner").
		Joins("INNER JOIN providers on quotations.provider_id = providers.id").
		Joins("INNER JOIN users  on quotations.user_id = users.id").
		Group("providers.id, providers.name, quotations.requirement_id, users.id, users.first_name, users.last_name, quotations.id, quotations.winner_level, quotations.Winner").
		Having("quotations.requirement_id = ?", request.RequirementID).
		Order("winner_level asc").
		Scan(&quotationResults).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Find quotations get prices ========= Quotation, QuotationDetail, Require
	quotationPricesResults := make([]winnerLevelResult, 0)
	if err := db.Table("quotations").
		Select("quotations.id, quotations.provider_id, quotations.requirement_id, sum(quotation_details.unit_price * requires.amount) as summation, quotations.suggest_winner").
		Joins("INNER JOIN quotation_details on quotations.id = quotation_details.quotation_id").
		Joins("INNER JOIN requires on quotation_details.require_id = requires.id").
		Group("quotations.provider_id, quotations.requirement_id, quotations.suggest_winner, quotations.id").
		Having("quotations.requirement_id = ?", request.RequirementID).
		Order("summation asc").
		Scan(&quotationPricesResults).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	responseQuotations := make([]quotationCustomResponse, 0)
	for _, qNames := range quotationResults {
		for _, v := range quotationPricesResults {
			if qNames.ID == v.ID {
				responseQuotations = append(responseQuotations, quotationCustomResponse{
					ID:            qNames.ID,
					ProviderID:    qNames.ProviderID,
					ProviderName:  qNames.ProviderName,
					UserID:        qNames.UserID,
					UserFirstName: qNames.UserFirstName,
					UserLastName:  qNames.UserLastName,
					RequirementID: qNames.RequirementID,
					Count:         qNames.Count,
					WinnerLevel:   qNames.WinnerLevel,
					Winner:        qNames.Winner,
					Summation:     v.Summation,
				})
			}
		}
	}
	total := len(responseQuotations)

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success: true,
		Data:    responseQuotations,
		Total:   uint(total),
	})
}

type purchaseOrder struct {
	Code          string  `json:"code"`
	Amount        string  `json:"amount"`
	ProviderID    uint    `json:"-"`
	RequirementID uint    `json:"-"`
	UnitMeasure   string  `json:"unit_measure"`
	Description   string  `json:"description"`
	UnitPrice     float32 `json:"unit_price"`
	Total         float32 `json:"total"`
}

type purchaseOrderResponse struct {
	PurchaseOrder []purchaseOrder    `json:"purchase_order"`
	Provider      models.Provider    `json:"provider"`
	Requirement   models.Requirement `json:"requirement"`
}

func PurchaseOrder(c echo.Context) error {
	// Get data request
	quotation := models.Quotation{}
	if err := c.Bind(&quotation); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	purchaseOrders := make([]purchaseOrder, 0)
	if err := db.Table("quotations").
		Select("quotations.id, quotations.provider_id, quotations.requirement_id, requires.amount, requires.unit_measure, products.name as description, quotation_details.unit_price, requires.amount * quotation_details.unit_price as total").
		Joins("INNER JOIN quotation_details on quotations.id = quotation_details.quotation_id").
		Joins("INNER JOIN requires on quotation_details.require_id = requires.id").
		Joins("INNER JOIN products on requires.product_id = products.id").
		Where("winner = true AND quotations.requirement_id = ?", quotation.RequirementID).
		Scan(&purchaseOrders).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	provider := models.Provider{}
	if err := db.First(&provider, purchaseOrders[0].ProviderID).Error; err != nil {
		return err
	}

	requirement := models.Requirement{}
	if err := db.First(&requirement, purchaseOrders[0].RequirementID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data: purchaseOrderResponse{
			PurchaseOrder: purchaseOrders,
			Provider:      provider,
			Requirement:   requirement,
		},
	})
}

type ctResponseRequire struct {
	ID          uint    `json:"id"`
	Amount      float32 `json:"amount"`
	Name        string  `json:"name"`
	UnitMeasure string  `json:"unit_measure"`
	Observation string  `json:"observation"`
}

// Quotations
type ctResponseQuotation struct {
	QuotationID uint    `json:"quotation_id"`
	UnitPrice   float32 `json:"unit_price"`
	Sequence    uint    `json:"sequence"`
}

// Providers
type ctResponseProvider struct {
	Name        string    `json:"name"`
	Manager     string    `json:"manager"`
	DeliverDate time.Time `json:"deliver_date"`
}

type comparativeTable struct {
	CTResponseQuotations []ctResponseQuotation `json:"ct_response_quotations"`
	CTResponseRequires   []ctResponseRequire   `json:"ct_response_requires"`
	CTResponseProviders  []ctResponseProvider  `json:"ct_response_providers"`
}

func ComparativeTable(c echo.Context) error {
	// Get data request
	requirement := models.Requirement{}
	if err := c.Bind(&requirement); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// -----------------------------------------------------------
	// Query Requires ------------------------------------------
	// -----------------------------------------------------------
	ctResponseRequires := make([]ctResponseRequire, 0)
	if err := db.Table("requires").
		Select("requires.id, products.name, requires.amount, requires.unit_measure").
		Joins("INNER JOIN products on requires.product_id = products.id").
		Where("requires.requirement_id = ?", requirement.ID).
		Order("requires.id asc").
		Scan(&ctResponseRequires).Error; err != nil {
		return err
	}

	// -----------------------------------------------------------
	// Query Quotations Quotation Detail  ------------------------
	// -----------------------------------------------------------
	ctResponseQuotations := make([]ctResponseQuotation, 0)
	if err := db.Table("providers").
		Select("quotations.id as quotation_id, quotation_details.unit_price").
		Joins("INNER JOIN quotations on providers.id = quotations.provider_id").
		Joins("INNER JOIN quotation_details on quotations.id = quotation_details.quotation_id").
		Where("quotations.requirement_id = ?", requirement.ID).
		Order("quotations.winner_level asc").
		Scan(&ctResponseQuotations).Error; err != nil {
		return err
	}

	// -----------------------------------------------------------
	// Query Providers  ------------------------------------------
	// -----------------------------------------------------------
	ctResponseProviders := make([]ctResponseProvider, 0)
	if err := db.Table("quotations").
		Select("providers.name, providers.manager, quotations.deliver_date").
		Joins("INNER JOIN providers on quotations.provider_id = providers.id").
		Where("quotations.requirement_id = ?", requirement.ID).
		Order("quotations.winner_level asc").
		Scan(&ctResponseProviders).Error; err != nil {
		return err
	}

	// Add column Sequence
	sequence := 1
	separate := len(ctResponseRequires)
	for i := 0; i < len(ctResponseQuotations); i++ {
		if i >= separate {
			sequence++
			separate = len(ctResponseRequires) * sequence
		}
		ctResponseQuotations[i].Sequence = uint(sequence)
	}

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: comparativeTable{
			CTResponseQuotations: ctResponseQuotations,
			CTResponseRequires:   ctResponseRequires,
			CTResponseProviders:  ctResponseProviders,
		},
	})
}

func CalculateWinnerLevelQuotation(requirementID uint) {
	// get connection
	db := config.GetConnection()
	defer db.Close()

	// CONSULT DATABASE
	quotationResults := make([]winnerLevelResult, 0)
	if err := db.Table("quotations").
		Select("quotations.id, quotations.provider_id, quotations.requirement_id, sum(quotation_details.unit_price * requires.amount) as summation, quotations.suggest_winner").
		Joins("INNER JOIN quotation_details on quotations.id = quotation_details.quotation_id").
		Joins("INNER JOIN requires on quotation_details.require_id = requires.id").
		Group("quotations.provider_id, quotations.requirement_id, quotations.suggest_winner, quotations.id").
		Having("quotations.requirement_id = ?", requirementID).
		Order("summation asc").
		Scan(&quotationResults).Error; err != nil {
		log.Panic(err)
	}

	// Update database
	for k, winnerQ := range quotationResults {
		quotation := models.Quotation{
			ID:          winnerQ.ID,
			WinnerLevel: uint(k) + 1,
		}
		db.Model(&quotation).Update(quotation)
	}
}

func CalculateWinnerByQuotation(requirementID uint) uint {
	return 0
}

// SetWinnerProvider se winner final provider in quotation
// RequestQuotation.ID == 0  -> Automatic calculate     // Optional
// RequestQuotation.ID != 0  -> Manual calculate        // Optional
// Requirement.ID                                       // Required
func SetWinnerQuotation(c echo.Context) error {
	// Get data request
	request := utilities.RequestQuotation{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validate if Manual or automatic calculation of the winner
	var WinnerID uint
	if request.ID == 0 {
		WinnerID = CalculateWinnerByQuotation(request.RequirementID) // Automatic calculate
	}
	WinnerID = request.ID

	quotation := models.Quotation{
		ID: WinnerID,
	}

	// Update all winners in false
	if err := db.Table("quotations").Where("requirement_id = ?", request.RequirementID).
		Updates(map[string]interface{}{"winner": false}).Error; err != nil {
		return err
	}

	// Update table quotation
	rows := db.Model(&quotation).UpdateColumn("winner", true).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", quotation.ID),
		})
	}

	// Change state requirement
	req := models.Requirement{
		ID:    request.RequirementID,
		State: "3",
	}
	rows = db.Model(&req).Update(req).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo cambiar de estado del requerimiento con el id = %d verfique que este requerimiento existe en la  base de datos", req.ID),
		})
	}

	// Return response success
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quotation.ID,
		Message: fmt.Sprintf("El ganador de la cotizacion con el id = %d se realizo exitosamente", quotation.ID),
	})
}

func GetQuotationByID(c echo.Context) error {
	// Get data request
	quotation := models.Quotation{}
	if err := c.Bind(&quotation); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&quotation, quotation.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quotation,
	})
}

func CreateQuotation(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	quotation := models.Quotation{}
	if err := c.Bind(&quotation); err != nil {
		return err
	}
	quotation.UserID = currentUser.ID
	quotation.EmissionDate = time.Now()

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Get Limit number quotations
	setting := models.Setting{}
	db.First(&setting)

	// Validate limit quotations
	var count uint
	if err := db.Model(&models.Quotation{}).Where("requirement_id = ?", quotation.RequirementID).Count(&count).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if count >= setting.Quotations {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: "A alcanzado el número maximo de cotizaciones",
		})
	}

	// Insert quotation in database
	if err := db.Create(&quotation).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Change state requirement
	req := models.Requirement{
		ID:    quotation.RequirementID,
		State: "1",
	}
	rows := db.Model(&req).Update(req).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo cambiar de estado del requerimiento con el id = %d verfique que este requerimiento existe en la  base de datos", req.ID),
		})
	}

	// Winner level calculate in database
	CalculateWinnerLevelQuotation(quotation.RequirementID)

	// Return response success
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quotation.ID,
		Message: fmt.Sprintf("La cotizacion %d se realizo exitosamente", quotation.ID),
	})
}

func UpdateQuotation(c echo.Context) error {
	// Get data request
	quotation := models.Quotation{}
	if err := c.Bind(&quotation); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Prepare data to UPDATE
	details := quotation.QuotationDetails
	onlyQuotation := quotation
	onlyQuotation.QuotationDetails = []models.QuotationDetail{}

	// Update quotation
	rows := db.Model(&onlyQuotation).Update(onlyQuotation).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", quotation.ID),
		})
	}

	// Update quotation details
	for _, qd := range details {
		if err := db.Model(&qd).UpdateColumn("unit_price", qd.UnitPrice).Error; err != nil {
			return err
		}
	}

	// Winner level calculate in database
	CalculateWinnerLevelQuotation(quotation.RequirementID)

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quotation.ID,
	})
}

func DeleteQuotation(c echo.Context) error {
	// Get data request
	quotation := models.Quotation{}
	if err := c.Bind(&quotation); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation quotation exist
	if db.First(&quotation).RecordNotFound() {
		return c.JSON(http.StatusCreated, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", quotation.ID),
		})
	}

	// Delete quotation in database
	if err := db.Delete(&quotation).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quotation.ID,
	})
}
