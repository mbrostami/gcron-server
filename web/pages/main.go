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

// GetMethods method available for this page
func (p *MainPage) GetMethods() []string {
	return []string{"GET"}
}

// Handler get page parameters
func (p *MainPage) Handler(method string, c *gin.Context) Response {
	var res Response
	res = gin.H{
		"message": "main pong",
	}
	c.JSON(200, res)
	return res
}
