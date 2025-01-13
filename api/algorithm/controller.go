package algorithm

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AlgorithmController(c *gin.Context) {
	results, err := GenerateRecommendation(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}
