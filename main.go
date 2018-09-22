package main

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/requirement/api"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize migration database
	migration()

	// COR
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE, echo.PUT},
	}))

	// Static Files =========================================================================
	static := e.Group("/static")
	static.Static("", "static")

	// API
	api.PublicApi(e)
	api.ProtectedApi(e)

	// Custom port
	port := os.Getenv("PORT")
	if port == "" {
		port = config.GetConfig().Server.Port
	}

	// Starting server echo
	e.Logger.Fatal(e.Start(":" + port))
}

// migration Init migration database
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

	// -------------------------------------------------------------
	// INSERT FIST DATA --------------------------------------------
	// -------------------------------------------------------------
	usr := models.User{}
	db.First(&usr)
	// hash password
	cc := sha256.Sum256([]byte("admin"))
	pwd := fmt.Sprintf("%x", cc)
	// create model
	user := models.User{
		UserName: "admin",
		Password: pwd,
		Profile:  "admin",
		Email:    "yoel.antezana@gmail.com",
	}
	// insert database
	if usr.ID == 0 {
		db.Create(&user)
	}

	// First Setting
	cg := models.Setting{}
	db.First(&cg)
	co := models.Setting{
		Item:             10,
		CompanyName:      "REQUIREMENT WEB",
		CompanyShortName: "RW",
		Quotations:       3,
		Logo:             "static/logo.png",
	}
	// Insert database
	if cg.ID == 0 {
		db.Create(&co)
	}
}
