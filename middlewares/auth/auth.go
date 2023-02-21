package auth

import (
	"log"

	"github.com/gin-gonic/gin"
)

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func TokenAuthMiddleware(apiToken string) gin.HandlerFunc {
	requiredToken := apiToken

	if requiredToken == "" {
		log.Fatal("Please set API token in server execution params")
	}

	return func(c *gin.Context) {
		token := c.Request.Header.Get("X-Auth-Token")

		if token == "" {
			_, err := c.Cookie("gin_cookie")
			if err != nil {
				respondWithError(c, 401, "API token required")
				return
			}
		}

		if token != requiredToken {
			respondWithError(c, 401, "Invalid API token")
			return
		}

		c.SetCookie("gin_cookie", "lol", 3600, "/", "localhost", false, true)

		c.Next()
	}
}
