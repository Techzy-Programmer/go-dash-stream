package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func diskReact(uriPrefix, streamDir string) gin.HandlerFunc {
	if _, err := os.Stat(streamDir); os.IsNotExist(err) {
		log.Fatalf("Stream data directory does not exist: %v", streamDir)
	}

	fileserver := http.FileServer(http.Dir(streamDir))

	if uriPrefix != "" {
		fileserver = http.StripPrefix(uriPrefix, fileserver)
	}

	return func(c *gin.Context) {
		fileserver.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func StartStreamServer(port string) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(diskReact("/", StreamPath))

	fmt.Printf("Starting server on port %v\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
