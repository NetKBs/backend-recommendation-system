package auth

import (
	"example/api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", LoginController)
		auth.POST("/register", RegisterController)

		// Protected Route. Example of how to protect a route (Delete later)
		auth.GET("/test", middleware.AuthMiddleware(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Acceso concedido"})
		})
	}
}
