package router

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
)

func TypescriptTranspilationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		extension := filepath.Ext(r.URL.Path)

		if extension == ".js" {
			filename := strings.ReplaceAll(r.URL.Path, ".js", ".ts")
			cmd := exec.Command("esbuild", "web/src/"+filename, "--bundle")
			output, err := cmd.Output()
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}

			w.Header().Add("Content-Type", "application/javascript")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string(output))
			return
		}

		next.ServeHTTP(w, r)
	})
}
