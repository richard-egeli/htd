package router

import "github.com/a-h/templ"

type HtdPage struct {
	Template func(...interface{}) templ.Component
}
