package router

type Method string

const (
	GET     Method = "GET"
	PUT     Method = "PUT"
	HEAD    Method = "HEAD"
	POST    Method = "POST"
	PATCH   Method = "PATCH"
	DELETE  Method = "DELETE"
	OPTIONS Method = "OPTIONS"
)
