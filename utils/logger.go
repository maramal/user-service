package utils

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

/**
 * Aagrega el archivo server.log al output predeterminado del servidor Gin
 */
func SetupLogOutput() {
	file, _ := os.Create("server.log")
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
}
