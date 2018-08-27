package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/controller"
	"github.com/paulantezana/requirement/utilities"
)

// PublicApi public routes
func PublicApi(e *echo.Echo) {
	pb := e.Group("/api/v1")

	pb.POST("/user/login", controller.Login)
	pb.POST("/user/forgot/search", controller.ForgotSearch)
	pb.POST("/user/forgot/validate", controller.ForgotValidate)
	pb.POST("/user/forgot/change", controller.ForgotChange)
}

// ProtectedApi protected api token jwt
func ProtectedApi(e *echo.Echo) {
	ar := e.Group("/api/v1")

	// Configure middleware with the custom claims type
	con := middleware.JWTConfig{
		Claims:     &utilities.Claim{},
		SigningKey: []byte(config.GetConfig().Server.Key),
	}
	ar.Use(middleware.JWTWithConfig(con))

	// Crud user
	ar.POST("/user/all", controller.GetUsers)
	ar.POST("/user/byid", controller.GetUserByID)
	ar.POST("/user", controller.CreateUser)
	ar.PUT("/user", controller.UpdateUser)
	ar.DELETE("/user", controller.DeleteUser)
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
}
