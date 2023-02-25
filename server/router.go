package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"tonic-file-access-server/config"

	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
)

type File struct {
	Name    string
	Size    int64
	ModTime time.Time
}

func NewRouter(apiToken string) *gin.Engine {
	dst := config.SetupDir()

	// // Disable Console Color, you don't need console color when writing the logs to file.
	// gin.DisableConsoleColor()

	// // Logging to a file.
	// f, _ := os.Create("gin.log")
	// gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()

	//router.Use(auth.TokenAuthMiddleware(apiToken))

	//load assets path
	router.StaticFile("/tonic.webp", "./assets/tonic.webp")
	router.StaticFile("/ginglass.mp4", "./assets/ginglass.mp4")
	//load templates
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/", func(c *gin.Context) {
		//query the /listdirectory endpoint to get the list of avaliable files
		//and pass it to the index.tmpl
		resp, _ := http.Get("http://localhost:5000/listdirectory")
		responseCode := resp.StatusCode
		if responseCode != 200 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to list the files",
			})

			return
		}
		//parse the json response
		jsonParsed, err := gabs.ParseJSONBuffer(resp.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable parse json",
			})
			return
		}
		files, _ := jsonParsed.Search("files").Children()
		filestruct := make([]File, 0)
		for _, child := range files {
			newName := child.Search("Name").Data().(string)
			newSize := child.Search("Size").Data().(float64)
			newModTime, _ := time.Parse("2006-01-02 15:04", child.Search("ModTime").Data().(string))
			filestruct = append(filestruct, File{Name: newName, Size: int64(newSize), ModTime: newModTime})
		}

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Home Page",
			"files": filestruct,
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

			fullFiles = append(fullFiles, File{Name: stats.Name(), Size: stats.Size(), ModTime: stats.ModTime()})
		}

		c.JSON(http.StatusOK, gin.H{"files": fullFiles})
	})

	return router
}
