package movie

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMoviesController(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	title := c.Query("title")
	genre := c.Query("genre")
	var response map[string]interface{}

	if title != "" || genre != "" {
		movies, err := GetMoviesByFilterRepository(title, genre, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response = gin.H{"page": page, "limit": limit, "movies": movies}

	} else {
		movies, err := GetMoviesRepository(page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response = gin.H{"page": page, "limit": limit, "movies": movies}
	}

	c.JSON(http.StatusOK, response)
}

func GetMovieByIdController(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	movie, err := GetMovieByIdRepository(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}
