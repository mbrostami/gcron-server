package web

import (
	// Import the gorilla/mux library we just installed

	"html/template"
	"strings"
	"time"

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
	addPage(r, pages.NewTaskPage(db))

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
	t := template.New("")
	t.Funcs(template.FuncMap{
		"byteToString": func(value []byte) template.HTML {
			return template.HTML(strings.Replace(string(value), "\n", "<br>", -1))
		},
	}).Funcs(template.FuncMap{
		"secondsToDate": func(value int64) template.HTML {
			unixTimeUTC := time.Unix(value, 0)
			return template.HTML(unixTimeUTC.Format(time.RFC3339))
		},
	})
	_, err := t.ParseGlob("web/static/*.tmpl")
	return t, err
}
