package main

import (
	"TerrainGenBackend/dla"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"net/http"
	"os"
)

func AllowCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func AllowCorsFile(fs http.Handler) http.HandlerFunc {
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
	Width      int
	Height     int
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

			img.SetRGBA(x, y, color.RGBA{R: 128, G: 128, B: 128, A: 255})
			// img.SetRGBA(x, y, color.RGBA{R: uint8(rand.Int()), G: uint8(rand.Int()), B: uint8(rand.Int()), A: 255})
		}
	}
	imgf, err := os.Create("textures/map.jpg")
	jpeg.Encode(imgf, img, &jpeg.Options{Quality: 100})
	g := dla.NewGrid(d.Width/9, d.Height/9, false)
	// g.RunDLACycles(25, 50)
	// g.PrintGrid()
	// g.UpscaleBy3()
	// g.RunDLACycles(200, 1000)
	// g.PrintGrid()
	// g.UpscaleBy3()
	// g.RunCrystalGrowth(0.9, 1)
	// g.PrintGrid()
	// g.CalculateEndDistance()
	// g.PrintGrid()
	// g.RunDLACycles(1000, 1000)
	// g.PrintGrid()
	f := g.ToFloatGrid(true)
	// f.BoxBlur(5, true)
	// f.BoxBlur(3, true)
	// f.BoxBlur(1, true)
	// f.BoxBlur(1, true)
	// f.BoxBlur(1, true)
	f = dla.NewFloatGrid(d.Width, d.Height)
	f.SimplexFill(5, 7.0)
	f.CircleFilter(10, 10.0)
	f.Normalize()
	// f.BoxBlur(1, false)
	// f.BoxBlur(1, false)
	// f.BoxBlur(1, false)
	// f.BoxBlur(1, false)
	// f.BoxBlur(1, true)
	// f.BoxBlur(3, true)
	// f.BoxBlur(1, true)
	// f.BoxBlur(1, true)
	// f.BoxBlur(1, true)
	// f.BoxBlur(1, true)
	var reply replyData
	reply.Width = d.Width
	reply.Height = d.Height
	reply.Heights = f.ExportHeights()
	reply.TextureURL = "http://localhost:8080/tex/map.jpg"

	j, _ := json.Marshal(reply)

	w.Header().Set("Content-Type", "text/plain")
	// w.Header().Set("Content-Type", "m")
	fmt.Fprintf(w, "%s", j)
}
