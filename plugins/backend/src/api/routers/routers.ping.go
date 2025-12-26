package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(_c *gin.Context) {
	_c.JSON(http.StatusOK, "pong")
}
