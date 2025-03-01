package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

// Simple development server for testing the landing pages
func main() {
	port := flag.Int("port", 8080, "port to serve on")
	dir := flag.String("dir", "./", "directory to serve from")
	flag.Parse()

	fs := http.FileServer(http.Dir(filepath.Clean(*dir)))
	http.Handle("/", fs)

	log.Printf("Starting server on :%d, serving from directory: %s\n", *port, *dir)
	log.Printf("Access at: http://localhost:%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
