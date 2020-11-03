package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var codes []string;

func main() {
	codes = append(codes, "abc", "promo", "avancadev", "dev")
	http.HandleFunc("/", get)
	http.ListenAndServe(":9093", nil)
}

func get(w http.ResponseWriter, r *http.Request)  {
	jsoResult, err := json.Marshal(codes)

	if err != nil {
		log.Fatal("Erro processing json")
	}

	fmt.Fprintf(w, string(jsoResult))
}