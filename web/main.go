package web

import (
	// Import the gorilla/mux library we just installed

	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mbrostami/gcron-server/configs"
	"github.com/mbrostami/gcron-server/db"
	"github.com/mbrostami/gcron-server/web/pages"
	pb "github.com/mbrostami/gcron/grpc"
	log "github.com/sirupsen/logrus"
)

// Listen start web server
func Listen(db db.DB, cfg configs.Config) {
	// gin.DefaultWriter = log.StandardLogger().Writer()
	r := gin.Default()
	t, _ := loadTemplate()
	r.SetHTMLTemplate(t)

	r.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))

	r.Use(authRequired)
	addPublicPage(r, pages.NewLoginPage(cfg.Auth.Username, cfg.Auth.Password))
	r.Use(static.Serve("/", static.LocalFile("web/static/public", false)))
	authorized := r.Group("/")
	authorized.Use(authRequired)
	{
		addRouterPage(authorized, pages.NewMainPage(db))
		addRouterPage(authorized, pages.NewTaskPage(db))
		addRouterPage(authorized, pages.NewLogoutPage())
	}

	r.Run("localhost:1401")
	log.Infof("Started listening on: %d (http)", 1401)
}

func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	c.Set("user", user)
	c.Next()
}

func addPublicPage(r *gin.Engine, page pages.Page) {
	addRouterPage(&r.RouterGroup, page)
}

func addRouterPage(r *gin.RouterGroup, page pages.Page) {
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
			return template.HTML(unixTimeUTC.Format("15:04:05"))
		},
	}).Funcs(template.FuncMap{
		"timestampToDate": func(value *timestamp.Timestamp) template.HTML {
			tme := time.Unix(value.Seconds, 0)
			return template.HTML(tme.Format("Aug 2 15:04:05"))
		},
	}).Funcs(template.FuncMap{
		"nanoToMili": func(value int32) template.HTML {
			res := fmt.Sprintf("%04f", float64(value)/float64(time.Millisecond))
			return template.HTML(res)
		},
	}).Funcs(template.FuncMap{
		"getDuration": func(task *pb.Task) string {
			durationSecond := task.EndTime.Seconds - task.StartTime.Seconds
			duration := fmt.Sprintf(
				"%d.%d",
				durationSecond,
				int32(task.EndTime.Nanos-task.StartTime.Nanos)/int32(time.Millisecond),
			)
			return duration
		},
	})
	_, err := t.ParseGlob("web/static/*.tmpl")
	return t, err
}
