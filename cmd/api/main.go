package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/GuilhermRodovalho/kalendar"
)

const (
	minYear = 1963
	maxYear = 9999
)

func main() {
	mux := http.NewServeMux()
	// Phase 1 endpoints
	mux.HandleFunc("GET /celebrations/{year}", handleCelebrations)
	mux.HandleFunc("GET /liturgical-year/{year}", handleLiturgicalYear)
	// Phase 2 endpoints
	mux.HandleFunc("GET /calendar/{year}/mobile-dates", handleMobileDates)
	mux.HandleFunc("GET /calendar/{year}/{month}/{day}", handleCalendarDay)
	mux.HandleFunc("GET /calendar/{year}", handleCalendar)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	log.Println("Server started on port 8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
	log.Println("Server stopped gracefully")
}

func parseYear(r *http.Request) (int, error) {
	year, err := strconv.Atoi(r.PathValue("year"))
	if err != nil {
		return 0, err
	}
	if year < minYear || year > maxYear {
		return 0, strconv.ErrRange
	}
	return year, nil
}

func isValidDate(year, month, day int) bool {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t.Year() == year && int(t.Month()) == month && t.Day() == day
}

func handleCelebrations(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	year, err := parseYear(r)
	if err != nil {
		http.Error(w, "invalid year (must be between 1963 and 9999)", http.StatusBadRequest)
		return
	}

	celebrations, err := kalendar.GetCelebrationsForYear(year)
	if err != nil {
		http.Error(w, "failed to load celebrations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(celebrations); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func handleCalendar(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	year, err := parseYear(r)
	if err != nil {
		http.Error(w, "invalid year (must be between 1963 and 9999)", http.StatusBadRequest)
		return
	}

	entries, err := kalendar.GetCalendar(year)
	if err != nil {
		http.Error(w, "failed to get calendar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func handleCalendarDay(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	year, err := parseYear(r)
	if err != nil {
		http.Error(w, "invalid year (must be between 1963 and 9999)", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(r.PathValue("month"))
	if err != nil || month < 1 || month > 12 {
		http.Error(w, "invalid month", http.StatusBadRequest)
		return
	}

	day, err := strconv.Atoi(r.PathValue("day"))
	if err != nil || day < 1 || day > 31 {
		http.Error(w, "invalid day", http.StatusBadRequest)
		return
	}

	if !isValidDate(year, month, day) {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	entry, err := kalendar.GetCalendarDay(year, kalendar.Month(month), day)
	if err != nil {
		http.Error(w, "failed to get calendar day", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func handleMobileDates(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	year, err := parseYear(r)
	if err != nil {
		http.Error(w, "invalid year (must be between 1963 and 9999)", http.StatusBadRequest)
		return
	}

	mobile := kalendar.GetMobileDates(year)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mobile); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func handleLiturgicalYear(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	year, err := parseYear(r)
	if err != nil {
		http.Error(w, "invalid year (must be between 1963 and 9999)", http.StatusBadRequest)
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
