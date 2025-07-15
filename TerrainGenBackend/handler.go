package main

import (
	"TerrainGenBackend/dla"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"math/rand"
	"net/http"
	"os"
)

func AllowCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func AllowCorsFile(fs http.Handler) http.HandlerFunc {
	fmt.Printf("Fuck meeee")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		fs.ServeHTTP(w, r)
	}

}

func Handler(w http.ResponseWriter, r *http.Request) {
	AllowCors(&w)
}

type jsonData struct {
	Width  int
	Height int
}

type replyData struct {
	Heights    []float64
	TextureURL string
}

func GenerateTerrain(w http.ResponseWriter, r *http.Request) {
	AllowCors(&w)
	bytes, err := io.ReadAll(r.Body)
	// fmt.Printf("json: %s\n", string(bytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var d jsonData
	err = json.Unmarshal(bytes, &d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	img := image.NewRGBA(image.Rect(0, 0, d.Width+1, d.Height+1))
	for x := 0; x < d.Width; x++ {
		for y := 0; y < d.Height; y++ {

			img.SetRGBA(x, y, color.RGBA{R: uint8(rand.Int()), G: uint8(rand.Int()), B: uint8(rand.Int()), A: 255})
		}
	}
	imgf, err := os.Create("textures/map.jpg")
	jpeg.Encode(imgf, img, &jpeg.Options{Quality: 100})
	g := dla.NewGrid(d.Width/9, d.Height/9, false)
	// g.RunDLACycles(25, 50)
	// g.PrintGrid()
	g.UpscaleBy3()
	// g.RunDLACycles(200, 1000)
	// g.PrintGrid()
	g.UpscaleBy3()
	// g.RunDLACycles(500, 1000)
	// g.PrintGrid()
	f := g.ToFloatGrid()
	// f.BoxBlur(1, true)
	// f.BoxBlur(1, true)
	// f.BoxBlur(1, true)
	// f.BoxBlur(1, true)
	// f.BoxBlur(1, true)
	// f.BoxBlur(1, true)
	// f.BoxBlur(1, true)
	var reply replyData
	reply.Heights = f.ExportHeights()
	reply.TextureURL = "http://localhost:8080/tex/map.jpg"

	j, _ := json.Marshal(reply)

	w.Header().Set("Content-Type", "text/plain")
	// w.Header().Set("Content-Type", "m")
	fmt.Fprintf(w, "%s", j)
}
