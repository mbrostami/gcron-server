package web

import (
	// Import the gorilla/mux library we just installed

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/mbrostami/gcron-server/db"
	"github.com/mbrostami/gcron-server/web/pages"
	log "github.com/sirupsen/logrus"
)

// Listen start web server
func Listen(db db.DB) {
	r := gin.Default()

	t, _ := loadTemplate()
	r.SetHTMLTemplate(t)

	addPage(r, pages.NewMainPage(db))
	addPage(r, pages.NewLoginPage(db))

	r.Run("localhost:1401")
	log.Infof("Started listening on: %d (http)", 1401)
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

func loadTemplate() (*template.Template, error) {
	template := template.New("")
	_, err := template.ParseGlob("web/static/*.tmpl")
	return template, err
}
