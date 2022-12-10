package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/qwuiemme/ellipsespace-server/docs"
	"github.com/qwuiemme/ellipsespace-server/internal/authorization"
	catalogueobject "github.com/qwuiemme/ellipsespace-server/internal/catalogue-object"
	serverstatus "github.com/qwuiemme/ellipsespace-server/internal/server-status"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}

// @Summary Add Object Catalogue
// @Security ApiKeyAuth
// @Tags MainAPI
// @Description Add a record of the object to the database.
// @Accept json
// @Produce json
// @Param Input body catalogueobject.CatalogueObject true "Object info"
// @Success 201 {object} serverstatus.StatusJson
// @Failure 400 {object} serverstatus.StatusJson
// @Failure 401
// @Router /api/add-object-catologue [post]
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

// @Summary Get Object Catalogue
// @Security ApiKeyAuth
// @Tags MainAPI
// @Description Returns an object record or null object with the passed name.
// @Accept json
// @Produce json
// @Param Input body catalogueobject.CatalogueObjectJsonGet true "Object name"
// @Success 200 {object} catalogueobject.CatalogueObject
// @Failure 400 {object} serverstatus.StatusJson
// @Failure 401
// @Failure 500 {object} serverstatus.StatusJson
// @Router /api/get-object-catologue [get]
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

// @Summary Get All Objects Catalogue
// @Security ApiKeyAuth
// @Tags MainAPI
// @Description Returns all object records in the database.
// @Produce json
// @Success 200 {object} []catalogueobject.CatalogueObject
// @Failure 401
// @Failure 500 {object} serverstatus.StatusJson
// @Router /api/get-all-object-catologue [get]
func getAllObjectCatalogueHandler(c *gin.Context) {
	slice, err := catalogueobject.GetAllFromDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, serverstatus.StatusJson{
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK, slice)
}

// @Summary Create Session
// @Tags Sessions
// @Description Writes a new session to the database and returns its Id.
// @Accept json
// @Produce json
// @Param Input body authorization.Session true "Session data"
// @Success 201 {number} int SessionID
// @Failure 400 {object} serverstatus.StatusJson
// @Failure 500 {object} serverstatus.StatusJson
// @Router /api/session/create [post]
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

	c.JSON(http.StatusCreated, obj.Id)
}

// @Summary Authorize in Session
// @Tags Sessions
// @Description Checks the entered data for correctness and returns the JWT token if the check is successful.
// @Accept json
// @Produce json
// @Param Input body authorization.SessionJsonGet true "Session data"
// @Success 200 {string} string JWT-token
// @Failure 400 {object} serverstatus.StatusJson
// @Failure 401 {object} serverstatus.StatusJson
// @Failure 500 {object} serverstatus.StatusJson
// @Router /api/session/auth [get]
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

// @Summary Update Session
// @Security ApiKeyAuth
// @Tags Sessions
// @Description Updates the session data with the specified Id.
// @Accept json
// @Produce json
// @Param Input body authorization.Session true "Session data"
// @Success 200 {object} serverstatus.StatusJson
// @Failure 400 {object} serverstatus.StatusJson
// @Failure 401
// @Failure 403 {object} serverstatus.StatusJson
// @Failure 500 {object} serverstatus.StatusJson
// @Router /api/session/update [put]
func updateSessionHandler(c *gin.Context) {
	obj, err := authorization.Unmarshal(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, serverstatus.StatusJson{
			Message: err.Error(),
		})

		return
	}

	if data, _ := authorization.ParseJWT(strings.Split(c.Request.Header.Get("Authorization"), " ")[1]); data.SessionName != obj.SessionName || (data.SessionName != obj.SessionName && data.AccessLevel != authorization.ADMIN_LEVEL) {
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

	if sb := authorization.ParseJWTFromHeader(c); sb.AccessLevel == authorization.ADMIN_LEVEL {
		err = obj.UpdateAll()
	} else {
		err = obj.Update()
	}

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
