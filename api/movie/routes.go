package movie

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	movies := r.Group("/movies")
	{
		movies.GET("/", GetMoviesController)
		movies.GET("/:id", GetMovieByIdController)
	}
}
