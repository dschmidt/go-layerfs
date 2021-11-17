package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	layerfs "github.com/dschmidt/go-layerfs/m"
)

func main() {
	port := "8090"

	upper, _ := filepath.Abs("examples/upper")
	lower, _ := filepath.Abs("examples/lower")

	fmt.Printf("upper: %s\n", upper)
	fmt.Printf("lower: %s\n", lower)

	layerFs := layerfs.New(os.DirFS(upper), os.DirFS(lower))

	fileServer := http.FileServer(http.FS(layerFs))
	fileServerHandler := http.StripPrefix("/files/", fileServer).ServeHTTP
	http.HandleFunc("/files/", fileServerHandler)

	log.Printf("Listening on :%s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		fmt.Printf("Could not start server on port :%s\n", port)
	}
}
