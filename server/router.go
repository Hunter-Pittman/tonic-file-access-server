package server

import (
	"net/http"
	"path/filepath"
	"strings"
	"tonic-file-access-server/config"
	"tonic-file-access-server/middlewares/auth"
	"tonic-file-access-server/middlewares/logger"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	dst := config.Setup()

	router := gin.Default()

	router.Use(logger.RequestID())
	router.Use(auth.TokenAuthMiddleware())

	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Home Page",
		})
	})

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
			return
		}

		if err := c.SaveUploadedFile(file, dst+file.Filename); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Your file has been successfully uploaded at /download-user-file/" + file.Filename})
	})

	router.GET("/download-user-file/:filename", func(c *gin.Context) {
		fileName := c.Param("filename")
		targetPath := filepath.Join(dst, fileName)

		if !strings.HasPrefix(filepath.Clean(targetPath), dst) {
			c.String(403, "Look like you attacking me")
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Type", "application/octet-stream")
		c.File(targetPath)
	})

	router.GET("/listdirectory", func(c *gin.Context) {
		files, err := filepath.Glob(dst + "*")
		var fileNames []string
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to list the files",
			})
			return
		}

		for _, path := range files {
			fileNames = append(fileNames, filepath.Base(path))
		}

		c.JSON(http.StatusOK, gin.H{"files": fileNames})
	})

	return router
}
