package main

import (
	"fmt"
	"log"
	_WSUtil "websocket/apps/websocket"
	"websocket/base"

	"github.com/gin-gonic/gin"

	"net/http"
)

func main() {
	// Generate a self-signed certificate
	certFile, keyFile, err := base.GenerateSelfSignedCert()
	if err != nil {
		log.Fatalf("Error in generating self-signed certificate: %v", err)
	}

	// fmt.Println("Certificate file:", certFile)
	// fmt.Println("Key file:", keyFile)

	base.GPortNo = "8443" // HTTPS port

	// Initialize Gin router
	r := gin.Default()

	// Serve static files (e.g., CSS, JS)
	r.Static("/static", "./static")

	// WebSocket server initialization
	wsserver := _WSUtil.NewWebSocketServer("*")

	// Define your routes
	r.GET(base.GBaseAPI+"/GetConnection", func(c *gin.Context) {
		wsserver.GetConnection(c.Writer, c.Request)
	})

	r.POST(base.GBaseAPI+"/SendData", func(c *gin.Context) {
		wsserver.GetData(c.Writer, c.Request)
	})

	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, base.GBaseAPI+"/")
	})

	// Serve a test chat HTML page
	r.GET(base.GBaseAPI+"/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chate.html", gin.H{})
	})

	// Load HTML templates
	r.LoadHTMLGlob("html/*")

	// Start the server with TLS
	log.Printf("WebSocket server is running at http://localhost:%s...", base.GPortNo)

	// if err := r.Run(":" + base.GPortNo); err != nil {
	// 	log.Fatalf("Server error: %v", err)
	// }
	if err := r.RunTLS(fmt.Sprintf(":%s", base.GPortNo), certFile, keyFile); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
