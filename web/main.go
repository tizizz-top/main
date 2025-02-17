package main

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/huxulm/main/web/ui"
)

func staticHandler() http.HandlerFunc {
	// Create a file system from the embedded files
	web, _ := fs.Sub(ui.Static, "out")
	h := http.FileServer(http.FS(web))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		path := r.URL.Path
		if strings.HasSuffix(path, "/shell") {
			r.URL.Path += ".html"
		}
		h.ServeHTTP(w, r)
	}
}

func main() {
	http.HandleFunc("/api/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello!"))
	})
	http.HandleFunc("/", staticHandler())
	http.ListenAndServe(":8080", nil)
}
