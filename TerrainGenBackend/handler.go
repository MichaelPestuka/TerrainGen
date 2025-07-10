package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// var coords [4]float32
	coords := [4]float32{0.1, 0.2, 0.3, 0.4}
	j, _ := json.Marshal(coords)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "json")
	fmt.Fprintf(w, "%s", j)
}
