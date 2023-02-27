package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tonic-file-access-server/config"
	"tonic-file-access-server/middlewares/auth"

	"github.com/gin-gonic/gin"
)

type File struct {
	Name    string
	Size    int64
	ModTime string
}

func NewRouter(apiToken string) *gin.Engine {
	dst := config.SetupDir()

	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()

	router.Use(auth.TokenAuthMiddleware(apiToken))

	//load assets path
	router.StaticFile("/tonic.webp", "./assets/tonic.webp")
	router.StaticFile("/ginglass.mp4", "./assets/ginglass.mp4")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	//load templates
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/", func(c *gin.Context) {

		//parse the json response
		files := listDirectory(dst)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Home Page",
			"files": files,
		})
	})

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["file[]"]

		for _, file := range files {
			if err := c.SaveUploadedFile(file, dst+file.Filename); err != nil {
				println(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to save the file",
				})
				return
			}
		}
		c.HTML(http.StatusOK, "upload.tmpl", gin.H{})
	})

	router.GET("/download/:filename", func(c *gin.Context) {
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
		var fullFiles []File
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to list the files",
			})
			return
		}

		for _, path := range files {
			stats, err := os.Stat(path)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to stat files",
				})
			}

			fullFiles = append(fullFiles, File{Name: stats.Name(), Size: (stats.Size() / 1024), ModTime: (stats.ModTime()).Format("2006-01-02 15:04")})
		}

		c.JSON(http.StatusOK, gin.H{"files": fullFiles})
	})

	return router
}

func listDirectory(dst string) []File {
	files, err := filepath.Glob(dst + "*")
	if err != nil {
		fmt.Println("%v", err)
	}
	var fullFiles []File

	for _, path := range files {
		stats, err := os.Stat(path)
		if err != nil {
			fmt.Println("%v", err)
		}
		fullFiles = append(fullFiles, File{Name: stats.Name(), Size: (stats.Size() / 1024), ModTime: (stats.ModTime()).Format("2006-01-02 15:04")})
	}

	return fullFiles

}
