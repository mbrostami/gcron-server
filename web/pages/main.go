package pages

import (
	"github.com/gin-gonic/gin"
)

// MainPage using Page interface
type MainPage struct {
	Template string
}

// NewMainPage creates new page
func NewMainPage() *MainPage {
	return &MainPage{}
}

// GetPath returns template path
func (p *MainPage) GetPath() string {
	return "web/static/main.html"
}

// GetRoute url endpoint
func (p *MainPage) GetRoute() string {
	return "/"
}

// GetMethod method available for this page
func (p *MainPage) GetMethod() string {
	return "GET"
}

// GetResponse get page parameters
func (p *MainPage) GetResponse() Response {
	var res Response
	res = gin.H{
		"message": "main pong",
	}
	return res
}
