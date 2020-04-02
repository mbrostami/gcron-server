package pages

import (
	"github.com/gin-gonic/gin"
)

// LoginPage using Page interface
type LoginPage struct {
	Template string
}

// NewLoginPage creates new page
func NewLoginPage() *LoginPage {
	return &LoginPage{}
}

// GetPath returns template path
func (p *LoginPage) GetPath() string {
	return "web/static/login.html"
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
	c.JSON(200, res)
	return res
}
