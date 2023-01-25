package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}

func aboutHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "about", nil)
}

func downloadHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "download", nil)
}

func techStackHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "tech-stack", nil)
}
