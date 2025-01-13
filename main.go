package main

import (
	"example/api/auth"
	"example/api/movie"
	"example/api/user"
	"example/config"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	config.ConnectDB()
}

func main() {
	defer config.SESSION.Close()

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	auth.RegisterRoutes(r)
	movie.RegisterRoutes(r)
	user.RegisterRoutes(r)

<<<<<<< HEAD
	r.GET("/algorithm/:user_id", algorithm.AlgorithmController)

	//algorithm.GenerateRecommendation("9aa2a501-4263-4049-af7d-9f13ad638b17")

=======
>>>>>>> eba1e67037d78fbf2ad4c3173d50172fc307b827
	r.Run()
}
