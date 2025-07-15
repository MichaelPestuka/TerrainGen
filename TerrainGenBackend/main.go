package main

import (
	// "TerrainGenBackend/httphandler"
	"net/http"
)

// "TerrainGenBackend/httphandler"
// "net/http"

func main() {

	textureServer := http.FileServer(http.Dir("./textures"))
	http.HandleFunc("/", Handler)
	http.Handle("GET /tex/", http.StripPrefix("/tex", AllowCorsFile(textureServer)))
	http.HandleFunc("QUERY /", GenerateTerrain)
	http.ListenAndServe(":8080", nil)
	// var g = dla.NewGrid(10, 10, true)
	// g.PrintGrid()
	// (&g).UpscaleBy3()
	// g.PrintGrid()
	// f := g.ToFloatGrid()
	// f.Print()
	// f.BoxBlur(1)
	// f.BoxBlur(1)
	// f.BoxBlur(1)
	// f.Print()
}
