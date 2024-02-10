package router

import (
	"net/http"
	"strings"
)

type HtdFileserver struct {
	Dir string
}

func (f *HtdFileserver) Create() {
	trimmed := strings.Trim(f.Dir, "/")
	slashed := "/" + trimmed + "/"

	fs := http.FileServer(http.Dir(trimmed))
	http.Handle(slashed, http.StripPrefix(slashed, fs))
}
