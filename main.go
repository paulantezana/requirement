package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/controller"
	"github.com/paulantezana/requirement/models"
	"github.com/paulantezana/requirement/utilities"
	"os"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	migration()
	//initial()
	//TestBorrame()

	// COR
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		//AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE, echo.PUT},
	}))

	// ======================================================================================
	// Static Files =========================================================================
	static := e.Group("/static")
	static.Static("", "static")

	// ======================================================================================
	// Public routes  =======================================================================
	pb := e.Group("/api/v1")

	pb.POST("/user/login", controller.Login)
	pb.POST("/user/forgout/serach", controller.ForgoutSearch)
	pb.POST("/user/forgout/validate", controller.ForgoutValidate)
	pb.POST("/user/forgout/change", controller.ForgoutChange)
	pb.POST("/user/register", controller.CreateUsuario)

	// ======================================================================================
	// CREATE GROUP API V1 SECRET ROUTES ====================================================
	ar := e.Group("/api/v1")

	// Configure middleware with the custom claims type
	con := middleware.JWTConfig{
		Claims:     &utilities.Claim{},
		SigningKey: []byte(config.GetConfig().Server.Key),
	}
	ar.Use(middleware.JWTWithConfig(con))

	// Crud user
	ar.POST("/user/all", controller.GetUsuarios)
	ar.POST("/user/byid", controller.GetUsuarioByID)
	ar.POST("/user", controller.CreateUsuario)
	ar.PUT("/user", controller.UpdateUsuario)
	ar.DELETE("/user", controller.DeleteUsuario)
	ar.POST("/user/upload/avatar", controller.UploadAvatarUser)
	ar.POST("/user/reset/password", controller.ResetPasswordUser)
	ar.POST("/user/change/password", controller.ChangePasswordUser)

	// Crud Product
	ar.POST("/product/all", controller.GetProducts)
	ar.POST("/product/byid", controller.GetProductByID)
	ar.POST("/product", controller.CreateProduct)
	ar.PUT("/product", controller.UpdateProduct)
	ar.DELETE("/product", controller.DeleteProduct)
	ar.POST("/product/search", controller.GetProductSearch)

	// Crud Provider
	ar.POST("/provider/all", controller.GetProviders)
	ar.POST("/provider/byid", controller.GetProviderByID)
	ar.POST("/provider", controller.CreateProvider)
	ar.PUT("/provider", controller.UpdateProvider)
	ar.DELETE("/provider", controller.DeleteProvider)
	ar.POST("/provider/search", controller.GetProviderSearch)

	// Crud Requirement
	ar.POST("/requirement/all", controller.GetRequirements)
	ar.POST("/requirement/byid", controller.GetRequirementByID)
	ar.POST("/requirement", controller.CreateRequirement)
	ar.PUT("/requirement", controller.UpdateRequirement)
	ar.DELETE("/requirement", controller.DeleteRequirement)
	ar.PUT("/requirement/set/rejected", controller.SetRejectedRequirement)
	ar.PUT("/requirement/set/closed", controller.SetClosedRequirement)

	// Crud Require
	ar.POST("/require/by/requirement", controller.GetRequireByRequirement)
	ar.POST("/require/by/quotation", controller.GetRequireByQuotation)
	ar.DELETE("/require", controller.DeleteRequire)

	// Crud Quotation
	ar.POST("/quotation/all", controller.GetQuotations)
	ar.POST("/quotation/byid", controller.GetQuotationByID)
	ar.POST("/quotation", controller.CreateQuotation)
	ar.PUT("/quotation", controller.UpdateQuotation)
	ar.DELETE("/quotation", controller.DeleteQuotation)
	ar.PUT("/quotation/set/winner", controller.SetWinnerQuotation)
	ar.POST("/quotation/comparativetable", controller.ComparativeTable)

	// Global settings
	ar.POST("/setting/global", controller.GetGlobalSettings)
	ar.GET("/setting", controller.GetSetting)
	ar.PUT("/setting", controller.UpdateSetting)

	// Statistic
	ar.POST("/statistic/top/provider/winners", controller.TopProviderWinner)
	ar.POST("/statistic/top/users", controller.TopUsers)
	ar.POST("/statistic/top/products", controller.TopProducts)
	ar.POST("/statistic/top/requirements", controller.TopRequirements)

	// Reporting EXCEL generate and Download
	ar.GET("/download/requirement/all", controller.ExportRequirementAll)

	// Custom port
	port := os.Getenv("PORT")
	if port == "" {
		port = config.GetConfig().Server.Port
	}

	// Starting server echo
	e.Logger.Fatal(e.Start(":" + port))
}

type ResultData struct {
	ID        uint   `json:"id"`
	UserName  string `json:"user_name"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
	UnitPrice uint   `json:"unit_price"`
}

type winnerLevelResult struct {
	ID            uint
	ProviderID    uint
	RequirementID uint
	Summation     float32
	Top           uint
}

func TestBorrame() {
	// get connection
	db := config.GetConnection()
	defer db.Close()

	// CONSULT DATABASE
	con := models.Setting{}
	db.First(&con)

	xlsx := excelize.NewFile()

	err := xlsx.AddPicture("Sheet1", "B2", "./static/logo.png", `{"x_scale": 0.5, "y_scale": 0.5}`)
	if err != nil {
		fmt.Println(err)
	}
	xlsx.SetCellValue("Sheet1", "A5", con.Company)
	xlsx.SetCellValue("Sheet1", "A6", con.City)
	xlsx.SetCellValue("Sheet1", "A8", "Requerimiento")

	//SET HEADER TABLE
	xlsx.SetCellValue("Sheet1", "A10", "Requerimiento")
	xlsx.SetCellValue("Sheet1", "B10", "Lugar")
	xlsx.SetCellValue("Sheet1", "C10", "Destino")
	xlsx.SetCellValue("Sheet1", "D10", "Fecha Emision")
	xlsx.SetCellValue("Sheet1", "E10", "Estado")

	// Get all requirements
	requirements := make([]models.Requirement, 0)
	db.Find(&requirements)

	currentRow := 11
	for k, rq := range requirements {
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", currentRow+k), rq.Name)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", currentRow+k), rq.Place)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", currentRow+k), rq.Destination)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", currentRow+k), rq.EmissionDate)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", currentRow+k), rq.State)
	}

	err = xlsx.SaveAs("./static/reports/requerimeinto.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("holaaaaaaaaaaaaa")
}

func migration() {
	db := config.GetConnection()
	defer db.Close()

	db.Debug().AutoMigrate(
		&models.User{},
		&models.Quotation{},
		&models.QuotationDetail{},
		&models.Product{},
		&models.Provider{},
		&models.Requirement{},
		&models.Require{},
		&models.Setting{},
	)
	db.Model(&models.Requirement{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Require{}).AddForeignKey("requirement_id", "requirements(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Require{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.QuotationDetail{}).AddForeignKey("require_id", "requires(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.QuotationDetail{}).AddForeignKey("quotation_id", "quotations(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Quotation{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Quotation{}).AddForeignKey("provider_id", "providers(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Quotation{}).AddForeignKey("requirement_id", "requirements(id)", "RESTRICT", "RESTRICT")
}
func initial() {
	db := config.GetConnection()
	defer db.Close()

	conf := models.Setting{
		Company:        "ABC Company",
		City:           "Sicuani - Cusco - Perú",
		Identification: "20563258412",
		Email:          "empresa@empresa.com",
		Item:           10,
		Quotations:     3,
	}

	if err := db.Create(&conf).Error; err != nil {
		log.Fatal(err)
	}
}
