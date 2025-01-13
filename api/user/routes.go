package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	user := r.Group("/users")
	{
		user.GET("", GetUsersController)
		user.GET("/:id", GetUserByIdController)
		user.PUT("/:id", UpdateUserController)
		user.DELETE("/:id", DeleteUserController)

		user.GET("/:id/history", GetUserHistoryController)
		user.GET("/:id/recommendations", GetUserRecommendationsController)
		user.POST("/:user_id/movies/:movie_id/watched", WatchMovieController)
	}
}
