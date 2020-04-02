package pages

import (
	"github.com/gin-gonic/gin"
	"github.com/mbrostami/gcron-server/db"
)

// LoginPage using Page interface
type LoginPage struct {
	db db.DB
}

// NewLoginPage creates new page
func NewLoginPage(db db.DB) *LoginPage {
	return &LoginPage{db: db}
}

// GetRoute url endpoint
func (p *LoginPage) GetRoute() string {
	return "/login"
}

// GetMethods method available for this page
func (p *LoginPage) GetMethods() []string {
	return []string{"GET", "POST"}
}

// Handler get page parameters
func (p *LoginPage) Handler(method string, c *gin.Context) Response {
	var res Response
	res = gin.H{
		"message": "login pong",
	}
	c.HTML(200, "login.tmpl", res)
	return res
}
