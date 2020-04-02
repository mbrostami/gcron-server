package pages

import "github.com/gin-gonic/gin"

type Page interface {
	Handler(method string, c *gin.Context) Response

	GetPath() string
	GetRoute() string
	GetMethods() []string
}

type Response interface {
}
