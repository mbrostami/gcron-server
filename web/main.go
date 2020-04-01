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
	r.GET(page.GetRoute(), func(c *gin.Context) {
		c.JSON(200, page.GetResponse())
	})
}
