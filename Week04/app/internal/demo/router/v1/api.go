package v1

import (
	"app/internal/demo/router/v1/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(app *gin.Engine) error {
	g := app.Group("/api")
	g.GET("/user", handler.User)
	return nil
}
