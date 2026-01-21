package main

import (
	"flag"
	"log"
	"net/http"
	"path"
	"strings"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dir := flag.String("dir", "web", "static directory")
	flag.Parse()

	fs := http.FileServer(http.Dir(*dir))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(path.Clean(r.URL.Path), ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")
		}
		fs.ServeHTTP(w, r)
	})

	log.Printf("serving %s on %s", *dir, *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal(err)
	}
}
