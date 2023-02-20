package server

import (
	"net/http"
	"path/filepath"
	"tonic-file-access-server/middlewares/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewRouter() *gin.Engine {
	// Env variables

	//dst := "C:\\test"

	router := gin.Default()

	router.Use(logger.RequestID())
	//router.Use(auth.TokenAuthMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")

		// The file cannot be received.
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
			return
		}

		// Retrieve file information
		extension := filepath.Ext(file.Filename)
		// Generate random file name for the new uploaded file so it doesn't override the old file with same name
		newFileName := uuid.New().String() + extension

		// The file is received, so let's save it
		if err := c.SaveUploadedFile(file, "C:\\test\\"+newFileName); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
			})
			return
		}

		// File saved successfully. Return proper result
		c.JSON(http.StatusOK, gin.H{"message": "Your file has been successfully uploaded."})
	})

	return router
}
