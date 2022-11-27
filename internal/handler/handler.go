package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	catalogueobject "github.com/qwuiemme/ellipsespace-server/internal/catalogue-object"
	servererror "github.com/qwuiemme/ellipsespace-server/internal/server-error"
)

func InitHandler() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("pages/*")

	router.GET("/", indexHandler)
	router.POST("/add-object-catologue", addObjectCatalogueHandler)
	router.POST("/get-object-catalogue", getObjectCatalogueHandler)

	return router
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}

func addObjectCatalogueHandler(c *gin.Context) {
	jsonByte, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, servererror.ErrorJson{
			Message: err.Error(),
		})

		return
	}

	var obj catalogueobject.CatalogueObject
	err = json.Unmarshal(jsonByte, &obj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, servererror.ErrorJson{
			Message: err.Error(),
		})

		return
	}

	err = obj.AddToDatabase()

	if err != nil {
		c.JSON(http.StatusBadRequest, servererror.ErrorJson{
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Complete!",
	})
}

// Необходимо переписать!
func getObjectCatalogueHandler(c *gin.Context) {
	jsonByte, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, servererror.ErrorJson{
			Message: err.Error(),
		})

		return
	}

	var obj catalogueobject.CatalogueObject
	err = json.Unmarshal(jsonByte, &obj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, servererror.ErrorJson{
			Message: err.Error(),
		})

		return
	}

	obj, err = catalogueobject.GetFromDatabase(obj.Name)

	if err != nil {
		c.JSON(http.StatusBadRequest, servererror.ErrorJson{
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, obj)
}
