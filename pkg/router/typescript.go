package router

import (
	"compress/gzip"
	"fmt"
	"log"
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
			cmd := exec.Command("esbuild", "web/src/"+filename, "--bundle", "--minify", "--format=esm", "--target=esnext")
			output, err := cmd.Output()
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}

			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Add("Content-Encoding", "gzip")
				w.Header().Add("Content-Type", "application/javascript")

				writer := gzip.NewWriter(w)
				defer writer.Close()

				_, err = writer.Write(output)
				if err != nil {
					log.Println("ERROR! ", err)
				}
				return
			}

			fmt.Fprint(w, string(output))
			w.WriteHeader(http.StatusOK)
		}

		next.ServeHTTP(w, r)
	})
}
