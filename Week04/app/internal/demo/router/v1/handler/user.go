package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func User(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
