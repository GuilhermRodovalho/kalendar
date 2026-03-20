package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/GuilhermRodovalho/kalendar"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /calendario-liturgico/{year}", handleLiturgicYear)

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleLiturgicYear(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	year, err := strconv.Atoi(r.PathValue("year"))
	if err != nil {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}

	ly := kalendar.LiturgicYearOf(year)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ly); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
