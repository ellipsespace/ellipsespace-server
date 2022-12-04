package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qwuiemme/ellipsespace-server/internal/authorization"
	catalogueobject "github.com/qwuiemme/ellipsespace-server/internal/catalogue-object"
	serverstatus "github.com/qwuiemme/ellipsespace-server/internal/server-status"
)

func InitHandler() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("pages/*")

	router.GET("/", indexHandler)
	api := router.Group("/api/")
	{
		sessions := api.Group("session/")
		{
			sessions.POST("create", createSessionHandler)
			sessions.GET("auth", authSessionHandler)
		}

		api.Use(authorization.AuthorizationRequired)

		api.GET("get-object-catalogue", getObjectCatalogueHandler)
		api.GET("get-all-object-catalogue", getAllObjectCatalogueHandler)
		sessions.PUT("update", updateSessionHandler)

		api.Use(authorization.AdminAccessLevelRequired)

		api.POST("add-object-catologue", addObjectCatalogueHandler)
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

	c.JSON(http.StatusCreated, serverstatus.StatusJson{
		Message: "Done.",
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
		c.JSON(http.StatusInternalServerError, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, obj)
}

func getAllObjectCatalogueHandler(c *gin.Context) {
	slice, err := catalogueobject.GetAllFromDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, serverstatus.StatusJson{
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK, slice)
}

func createSessionHandler(c *gin.Context) {
	obj, err := authorization.Unmarshal(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	err = obj.SetPassword(obj.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	err = obj.AddToDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	*obj, err = authorization.GetSession(obj.SessionName)

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": obj.Id,
	})
}

func authSessionHandler(c *gin.Context) {
	get, err := authorization.UnmarshalJsonGet(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	obj, err := authorization.GetSession(get.SessionName)

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	if obj.ComparePassword(get.Password) {
		token, err := authorization.GenerateJWT(obj.SessionBase)

		if err != nil {
			c.JSON(http.StatusInternalServerError, serverstatus.StatusJson{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, token)
	} else {
		c.JSON(http.StatusUnauthorized, serverstatus.StatusJson{
			Message: "The session name or password entered is incorrect.",
		})
	}
}

func updateSessionHandler(c *gin.Context) {
	obj, err := authorization.Unmarshal(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	if data, _ := authorization.ParseJWT(strings.Split(c.Request.Header.Get("Authorization"), " ")[1]); data.SessionName != obj.SessionName || (data.SessionName != obj.SessionName && data.AccessLevel != 1) {
		c.JSON(http.StatusForbidden, serverstatus.StatusJson{
			Message: "You are trying to edit someone else's session data without the necessary rights.",
		})

		return
	}

	err = obj.SetPassword(obj.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	err = obj.Update()

	if err != nil {
		c.JSON(http.StatusInternalServerError, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, serverstatus.StatusJson{
		Message: "Done.",
	})
}
