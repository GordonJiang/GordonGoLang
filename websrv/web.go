package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	path = flag.String("path", "", "Path to host the simple web server")
)

func main() {
	flag.Parse()
	if len(*path) < 2 {
		flag.Usage()
		return
	}
	log.Println("Starting web server on path: ", *path)

	log.Fatal(http.ListenAndServe(":8081", http.FileServer(http.Dir(*path))))
}
