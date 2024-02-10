package router

type HtdMethod string

const (
	GET     HtdMethod = "GET"
	PUT     HtdMethod = "PUT"
	HEAD    HtdMethod = "HEAD"
	POST    HtdMethod = "POST"
	PATCH   HtdMethod = "PATCH"
	DELETE  HtdMethod = "DELETE"
	OPTIONS HtdMethod = "OPTIONS"
)
