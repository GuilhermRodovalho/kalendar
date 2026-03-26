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
	mux.HandleFunc("GET /santos", handleSaints)
	mux.HandleFunc("GET /santos/{year}", handleSaintsByYear)

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

	ly, err := kalendar.GetLiturgicYearWithCelebrations(year)
	if err != nil {
		http.Error(w, "failed to get liturgical year", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ly); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func handleSaints(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	saints, err := kalendar.GetAllSaints()
	if err != nil {
		http.Error(w, "failed to load saints", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(saints); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func handleSaintsByYear(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	year, err := strconv.Atoi(r.PathValue("year"))
	if err != nil {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}

	saints, err := kalendar.GetSaintsForYear(year)
	if err != nil {
		http.Error(w, "failed to load saints", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(saints); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
