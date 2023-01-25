package handler

import (
	"github.com/gin-gonic/gin"
	_ "github.com/qwuiemme/ellipsespace-server/docs"
	"github.com/qwuiemme/ellipsespace-server/internal/authorization"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitHandler() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("pages/*")
	router.Static("/static", "static/")

	// Site handlers
	router.GET("/", indexHandler)
	router.GET("/tech", techStackHandler)
	router.GET("/about", aboutHandler)
	router.GET("/download", downloadHandler)

	// API handlers
	api := router.Group("/api/")
	{
		catalogue := api.Group("catalogue/")
		sessions := api.Group("session/")
		{
			sessions.POST("create", createSessionHandler)
			sessions.GET("auth", authSessionHandler)
		}

		api.Use(authorization.AuthorizationRequired)
		sessions.Use(authorization.AuthorizationRequired)

		catalogue.GET("get", getObjectCatalogueHandler)
		catalogue.GET("all", getAllObjectCatalogueHandler)
		sessions.PUT("update", updateSessionHandler)
		sessions.GET("info", idSessionHandler)

		api.Use(authorization.AdminAccessLevelRequired)
		sessions.Use(authorization.AdminAccessLevelRequired)

		catalogue.POST("add", addObjectCatalogueHandler)
		catalogue.PUT("update", updateObjectCatalogueHandler)
		catalogue.DELETE("delete", deleteObjectCatalogueHandler)
		sessions.DELETE("delete", deleteSessionHandler)
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
