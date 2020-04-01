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

// GetMethod method available for this page
func (p *LoginPage) GetMethod() string {
	return "GET"
}

// GetResponse get page parameters
func (p *LoginPage) GetResponse() Response {
	var res Response
	res = gin.H{
		"message": "login pong",
	}
	return res
}
