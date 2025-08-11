package main

import (
	"TerrainGenBackend/dla"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"time"
)

func DeleteTexture(filename string, delay int) {
	time.Sleep(time.Duration(delay) * time.Second)
	os.Remove(filename)
}

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
	Width    int
	Height   int
	Sealevel float64
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

	// Generate heightmap
	g := dla.NewGrid(d.Width, d.Height, false)
	g.SimplexFill(5, 7.0)
	g.CircleFilter(0.05, 10.0)
	g.Normalize()
	g.DrawOcean(d.Sealevel)
	// g.FindShallows(3)
	g.FillDepressions()
	g.OceanSloping(0.005)

	// Generate texture
	img := g.TerrainTypeTexture()
	imageTag := rand.Int() % 10000000
	imageFileName := "map" + strconv.Itoa(imageTag) + ".png"
	imgf, err := os.Create("textures/" + imageFileName)
	go DeleteTexture("textures/"+imageFileName, 15) // Delete texture after some time
	png.Encode(imgf, img)

	// Send reply
	var reply replyData
	reply.Width = d.Width
	reply.Height = d.Height
	reply.Heights = g.ExportHeights()
	reply.TextureURL = "http://localhost:8080/tex/" + imageFileName

	j, _ := json.Marshal(reply)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s", j)
}
