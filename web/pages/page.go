package pages

type Page interface {
	GetResponse() Response
	GetPath() string
	GetRoute() string
	GetMethod() string
}

type Response interface {
}
