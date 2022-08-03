package routers

// import "github.com/gin-gonic/gin"

import (
	"GIK_Web/src/middleware"
	"GIK_Web/src/routers/admin"
	"GIK_Web/src/routers/analytics"
	"GIK_Web/src/routers/auth"
	"GIK_Web/src/routers/client"
	"GIK_Web/src/routers/info"
	"GIK_Web/src/routers/invoice"
	"GIK_Web/src/routers/items"
	"GIK_Web/src/routers/items_temp"
	"GIK_Web/src/routers/location"
	"GIK_Web/src/routers/logs"
	"GIK_Web/src/routers/qr"
	"GIK_Web/src/routers/settings"
	"GIK_Web/src/routers/status"
	"GIK_Web/src/routers/tags"
	"GIK_Web/src/routers/transaction"
	"GIK_Web/utils"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	r.Use(middleware.CORSMiddleware())

	// r.Use(middleware.AuthMiddleware())

	r.GET("/ping", status.Ping)

	analyticsApi := r.Group("/analytics")
	{
		analyticsApi.Use(middleware.AuthMiddleware())
		analyticsApi.Use(middleware.AdvancedLoggingMiddleware())
		analyticsApi.GET("/transaction", analytics.GraphTransactions)
		analyticsApi.GET("/transaction/total", analytics.GraphTotalTransactions)
		analyticsApi.GET("/attention", analytics.AttentionRequired)
		analyticsApi.GET("/activity", analytics.GetRecentActivity)
		analyticsApi.GET("/trending", analytics.GetTrendingItems)
	}

	authApi := r.Group("/auth")
	{
		authApi.POST("/login", auth.Login)
		authApi.POST("/prelogin", auth.CheckPasswordForLogin)
		authApi.GET("/tfa", auth.CheckTfaStatusBeforeLogin)
		authApi.POST("/register", auth.Register)
		authApi.GET("/first_admin", auth.CreateFirstAdmin)
		authApi.GET("/scode", auth.GetSignupCodeInfo)
		authApi.GET("/status", middleware.AuthMiddleware(), auth.CheckAuthStatus)
		authApi.GET("/logout", auth.Logout)
	}

	itemsApi := r.Group("/items")
	{
		itemsApi.Use(middleware.AuthMiddleware())
		itemsApi.Use(middleware.AdvancedLoggingMiddleware())
		itemsApi.GET("/list", items.ListItem)
		itemsApi.GET("/export", items.ExportItems)
		itemsApi.GET("/lookup", items.LookupItem)
		itemsApi.POST("/bulklookup", items.GetIdsInBulk)
		itemsApi.PUT("/add", items.AddItem)
		itemsApi.DELETE("/delete", items.DeleteLocation)
		itemsApi.PATCH("/update", items.UpdateItem)
		itemsApi.GET("/suggest", items.GetAutoSuggest)
		stockApi := itemsApi.Group("/stock")
		{
			stockApi.PUT("/add", items.AddStock)
			stockApi.PUT("/remove", items.RemoveStock)
			stockApi.PATCH("/jump", items.JumpStock)
		}

		itemsApi.GET("/product_id", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"message": "random product id generated",
				"data":    utils.GenerateProductId(),
			})
		})
	}

	items2Api := r.Group("/itemstemp")
	{
		items2Api.Use(middleware.AuthMiddleware())
		items2Api.Use(middleware.AdvancedLoggingMiddleware())
		items2Api.GET("/list", items_temp.ListItem)
		items2Api.PUT("/add", items_temp.AddItem)
		items2Api.PATCH("/update", items_temp.UpdateItem)
		items2Api.DELETE("/delete", items_temp.DeleteLocation)
	}

	tagsApi := r.Group("/tags")
	{
		tagsApi.Use(middleware.AuthMiddleware())
		tagsApi.Use(middleware.AdvancedLoggingMiddleware())
		tagsApi.GET("/list", tags.ListTags)
		tagsApi.PUT("/add", tags.AddTags)
		tagsApi.PATCH("/update", tags.UpdateTags)
		tagsApi.DELETE("/delete", tags.DeleteTags)

	}

	clientsApi := r.Group("/client")
	{
		clientsApi.Use(middleware.AuthMiddleware())
		clientsApi.Use(middleware.AdvancedLoggingMiddleware())
		clientsApi.GET("/list", client.ListClient)
		// clientsApi.GET("/lookup", client.LookupLocation)
		clientsApi.PUT("/add", client.AddClient)
		clientsApi.DELETE("/delete", client.DeleteClient)
		clientsApi.PATCH("/update", client.UpdateClient)
	}

	locationsApi := r.Group("/location")
	{
		locationsApi.Use(middleware.AuthMiddleware())
		locationsApi.Use(middleware.AdvancedLoggingMiddleware())
		locationsApi.GET("/list", location.ListLocation)
		locationsApi.GET("/lookup", location.LookupLocation)
		locationsApi.PUT("/add", location.AddLocation)
		locationsApi.DELETE("/delete", location.DeleteLocation)
		locationsApi.PATCH("/update", location.UpdateLocation)
		locationsApi.GET("/scan", location.GetScannedData)
	}

	transactionApi := r.Group("/transaction")
	{
		transactionApi.Use(middleware.AuthMiddleware())
		transactionApi.Use(middleware.AdvancedLoggingMiddleware())
		transactionApi.GET("/list", transaction.ListTransactions)
		// transactionApi.GET("/listItem", transaction.ListTransactionItem)
		// transactionApi.PATCH("/updateItem", transaction.UpdateTransactionItem)
		transactionApi.PUT("/add", transaction.AddTransaction)
		// transactionApi.PUT("/addItem", transaction.AddTransactionItem)
		transactionApi.DELETE("/delete", transaction.DeleteTransaction)
		// transactionApi.DELETE("/deleteItem", transaction.DeleteTransactionItem)
		transactionApi.GET("/items", transaction.GetTransactionItems)
	}

	logsApi := r.Group("/logs")
	{
		logsApi.Use(middleware.AuthMiddleware())
		logsApi.Use(middleware.AdvancedLoggingMiddleware())
		logsApi.GET("/advanced", logs.GetAdvancedLogs)
		logsApi.GET("/simple", logs.GetSimpleLogs)
	}

	adminApi := r.Group("/admin")
	{
		adminApi.Use(middleware.AuthMiddleware())
		adminApi.Use(middleware.AdvancedLoggingMiddleware())
		adminApi.Use(middleware.AdminMiddleware())
		adminApi.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"message": "Admin verified",
			})
		})
		adminApi.GET("/scodes", admin.GetSignupCodes)
		adminApi.GET("/scode", admin.CreateSignupCode)
		adminApi.PATCH("/scode/toggle", admin.ToggleSignupCode)
		adminApi.DELETE("/scode/delete", admin.DeleteSignupCodes)
		adminApi.GET("/lists", admin.ListAdminsAndNonAdmins)
		adminApi.GET("/users", admin.ListUsers)
		adminApi.PATCH("/user/toggle", admin.ToggleUser)
		adminApi.DELETE("/user/delete", admin.DeleteUser)
		adminApi.PATCH("/admins", admin.EditAdmins)
	}

	settingsApi := r.Group("/settings")
	{
		settingsApi.Use(middleware.AuthMiddleware())
		settingsApi.Use(middleware.AdvancedLoggingMiddleware())
		settingsApi.PATCH("/password", settings.ChangePassword)
		tfaApi := settingsApi.Group("/tfa")
		{
			tfaApi.GET("/status", settings.GetTfaStatus)
			tfaApi.GET("/generate", settings.GenerateTwoFactorSecret)
			tfaApi.PATCH("/setup", settings.ValidateAndSetupTwoFactor)
		}
	}

	qrApi := r.Group("/qr")
	{
		qrApi.Use(middleware.AuthMiddleware())
		qrApi.Use(middleware.AdvancedLoggingMiddleware())
		qrApi.GET("/codes", qr.GetQRCodes)
	}

	invoiceApi := r.Group("/invoice")
	{
		invoiceApi.Use(middleware.AuthMiddleware())
		invoiceApi.Use(middleware.AdvancedLoggingMiddleware())
		invoiceApi.POST("/generate", invoice.GetInvoice)
	}

	infoApi := r.Group("/info")
	{
		infoApi.Use(middleware.AuthMiddleware())

		infoApi.GET("/username", info.GetUsername)
		infoApi.GET("/currentusername", info.GetCurrentUsername)
		infoApi.GET("/client", info.GetClientInfo)
	}

	return r
}
