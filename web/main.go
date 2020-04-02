package web

import (
	// Import the gorilla/mux library we just installed

	"github.com/gin-gonic/gin"
	"github.com/mbrostami/gcron-server/web/pages"
	log "github.com/sirupsen/logrus"
)

// Listen start web server
func Listen() {
	r := gin.Default()

	addPage(r, pages.NewMainPage())
	addPage(r, pages.NewLoginPage())

	log.Infof("Started listening on: %d (http)", 1401)
	r.Run("localhost:1401")
}

func addPage(r *gin.Engine, page pages.Page) {
	for _, method := range page.GetMethods() {
		if method == "GET" {
			r.GET(page.GetRoute(), func(c *gin.Context) {
				page.Handler("GET", c)
			})
		} else if method == "POST" {
			r.POST(page.GetRoute(), func(c *gin.Context) {
				page.Handler("POST", c)
			})
		}
	}
}
