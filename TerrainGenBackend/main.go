package main

import (
	"net/http"
	"os"
)

func main() {

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	textureServer := http.FileServer(http.Dir("./textures"))
	http.HandleFunc("/", Handler)
	http.Handle("GET /tex/", http.StripPrefix("/tex", AllowCorsFile(textureServer)))
	http.HandleFunc("QUERY /", GenerateTerrain)
	http.ListenAndServe(":"+httpPort, nil)
}
