package server

import (
	"fmt"
	"log"
	"net/http"
	"tonic-file-access-server/middlewares/auth"
	"tonic-file-access-server/middlewares/logger"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	// Env variables

	dst := "C:\\Users\\hunte\\Documents\\repos\\tonic-file-access-server\\uploads"

	router := gin.Default()

	router.Use(logger.RequestID())
	router.Use(auth.TokenAuthMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, err := c.MultipartForm()

		if err != nil {
			log.Fatal(err)
		}

		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	router.POST("/uploadsingle", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	return router
}
