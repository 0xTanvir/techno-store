package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func health(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "ok"})
}
