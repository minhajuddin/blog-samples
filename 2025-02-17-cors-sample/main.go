package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// TODO: Move to url helpers
	stripTrailingSlash := func(s string) string {
		url, err := url.Parse(s)
		if err != nil {
			return s
		}

		return fmt.Sprintf("%s://%s", url.Scheme, url.Host)
	}

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     Map([]string{"http://localhost:8000/"}, stripTrailingSlash),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           time.Duration(12) * time.Hour,
	})

	// Set up CORS
	r.Use(corsMiddleware)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "http://localhost:8000/")
	})

	rg := r.Group("/api/v1/")

	rg.POST("/login", func(c *gin.Context) {
		c.SetCookie("_auth", "awesome.jwt.token", 3600, "", "", false, true)

		c.JSON(http.StatusOK, gin.H{"message": "Logged in via cookie"})
	})

	rg.GET("/me", func(c *gin.Context) {
		jwt, err := c.Cookie("_auth")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Cookie not found"})
			return
		}

		if jwt != "awesome.jwt.token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": "minhajuddin", "id": 1})
	})

	// Start the server
	err := r.Run(":8001")
	if err != nil {
		panic(err)
	}
}

// TODO: Move to an enum package
// Map is a generic function to map a slice of values to another slice of values
func Map[T1 any, T2 any](s []T1, f func(T1) T2) []T2 {
	mapped := make([]T2, len(s))
	for i, v := range s {
		mapped[i] = f(v)
	}
	return mapped
}
