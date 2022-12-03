package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	catalogueobject "github.com/qwuiemme/ellipsespace-server/internal/catalogue-object"
	serverstatus "github.com/qwuiemme/ellipsespace-server/internal/server-status"
)

func InitHandler() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("pages/*")

	router.GET("/", indexHandler)
	api := router.Group("/api/")
	{
		api.POST("add-object-catologue", addObjectCatalogueHandler)
		api.POST("get-object-catalogue", getObjectCatalogueHandler)
		api.POST("get-all-object-catalogue", getAllObjectCatalogue)
	}

	return router
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}

func addObjectCatalogueHandler(c *gin.Context) {
	obj, err := catalogueobject.Unmarshal(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	err = obj.AddToDatabase()

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, serverstatus.StatusJson{
		Message: "Complete!",
	})
}

func getObjectCatalogueHandler(c *gin.Context) {
	name, err := catalogueobject.UnmarshalJsonGet(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	obj, err := catalogueobject.GetFromDatabase(name.Name)

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, obj)
}

func getAllObjectCatalogue(c *gin.Context) {
	slice, err := catalogueobject.GetAllFromDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, serverstatus.StatusJson{
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK, slice)
}
