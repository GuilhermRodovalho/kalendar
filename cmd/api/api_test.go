package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GuilhermRodovalho/kalendar"
)

func TestHandleLiturgicalYear200(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /liturgical-year/{year}", handleLiturgicalYear)

	req := httptest.NewRequest("GET", "/liturgical-year/2026", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp kalendar.LiturgicSeasonsWithCelebrations
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(resp.Celebrations) == 0 {
		t.Error("celebrations should not be empty")
	}

	if resp.LiturgicSeasons.Advent.Start.Year() == 0 {
		t.Error("Advent start should be set")
	}
}

func TestHandleLiturgicalYear400(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /liturgical-year/{year}", handleLiturgicalYear)

	req := httptest.NewRequest("GET", "/liturgical-year/abc", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleCelebrations200(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /celebrations/{year}", handleCelebrations)

	req := httptest.NewRequest("GET", "/celebrations/2025", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var celebrations []kalendar.Celebration
	if err := json.Unmarshal(w.Body.Bytes(), &celebrations); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(celebrations) == 0 {
		t.Error("celebrations should not be empty")
	}

	for _, c := range celebrations {
		if c.Name == "" {
			t.Error("celebration name should not be empty")
		}
	}
}

func TestHandleCelebrations400(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /celebrations/{year}", handleCelebrations)

	req := httptest.NewRequest("GET", "/celebrations/invalid", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestLiturgicalYearResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /liturgical-year/{year}", handleLiturgicalYear)

	req := httptest.NewRequest("GET", "/liturgical-year/2026", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Content-Type = %v, want application/json", contentType)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if _, ok := resp["liturgical_seasons"]; !ok {
		t.Error("response should contain 'liturgical_seasons'")
	}
	if _, ok := resp["celebrations"]; !ok {
		t.Error("response should contain 'celebrations'")
	}
	if _, ok := resp["mobile_dates"]; ok {
		t.Error("response should NOT contain 'mobile_dates'")
	}
}

func TestCelebrationsNotEmpty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /liturgical-year/{year}", handleLiturgicalYear)

	req := httptest.NewRequest("GET", "/liturgical-year/2026", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	var resp kalendar.LiturgicSeasonsWithCelebrations
	json.Unmarshal(w.Body.Bytes(), &resp)

	if len(resp.Celebrations) < 100 {
		t.Errorf("expected at least 100 celebrations, got %d", len(resp.Celebrations))
	}
}

func TestHandleLiturgicalYearLeapYear(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /liturgical-year/{year}", handleLiturgicalYear)

	req := httptest.NewRequest("GET", "/liturgical-year/2024", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp kalendar.LiturgicSeasonsWithCelebrations
	json.Unmarshal(w.Body.Bytes(), &resp)

	if len(resp.Celebrations) < 100 {
		t.Errorf("expected at least 100 celebrations, got %d", len(resp.Celebrations))
	}
}

func TestHandleLiturgicalYearVariousYears(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /liturgical-year/{year}", handleLiturgicalYear)

	tests := []struct {
		year string
	}{
		{"2025"},
		{"2026"},
		{"2027"},
	}

	for _, tt := range tests {
		t.Run(tt.year, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/liturgical-year/"+tt.year, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
			}

			var resp kalendar.LiturgicSeasonsWithCelebrations
			json.Unmarshal(w.Body.Bytes(), &resp)

			if len(resp.Celebrations) < 100 {
				t.Errorf("expected at least 100 celebrations, got %d", len(resp.Celebrations))
			}
		})
	}
}

func TestCelebrationsGradeDistribution(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /celebrations/{year}", handleCelebrations)

	req := httptest.NewRequest("GET", "/celebrations/2026", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	var celebrations []kalendar.Celebration
	json.Unmarshal(w.Body.Bytes(), &celebrations)

	grades := make(map[kalendar.CelebrationGrade]int)
	for _, c := range celebrations {
		grades[c.Grade]++
	}

	if grades[kalendar.GradeSolemnity] == 0 {
		t.Error("should have at least one Solemnity")
	}
	if grades[kalendar.GradeFeast] == 0 {
		t.Error("should have at least one Feast")
	}
	if grades[kalendar.GradeMemorial] == 0 {
		t.Error("should have at least one Memorial")
	}
}

func TestCelebrationsColorDistribution(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /celebrations/{year}", handleCelebrations)

	req := httptest.NewRequest("GET", "/celebrations/2026", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	var celebrations []kalendar.Celebration
	json.Unmarshal(w.Body.Bytes(), &celebrations)

	colors := make(map[kalendar.LiturgicalColor]int)
	for _, c := range celebrations {
		colors[c.Color]++
	}

	if colors[kalendar.White] == 0 {
		t.Error("should have celebrations with white color")
	}
	if colors[kalendar.Red] == 0 {
		t.Error("should have celebrations with red color")
	}
}

func TestCelebrationsIncludesMobile(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /celebrations/{year}", handleCelebrations)

	req := httptest.NewRequest("GET", "/celebrations/2026", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	var celebrations []kalendar.Celebration
	json.Unmarshal(w.Body.Bytes(), &celebrations)

	mobileNames := map[string]bool{
		"Santos Pedro e Paulo, apóstolos":            false,
		"Todos os Santos":                            false,
		"Epifania do Senhor":                         false,
		"Nosso Senhor Jesus Cristo, Rei do Universo": false,
	}

	for _, c := range celebrations {
		if _, ok := mobileNames[c.Name]; ok {
			mobileNames[c.Name] = true
			if !c.IsMovable {
				t.Errorf("mobile celebration %q should have IsMovable=true", c.Name)
			}
		}
	}

	for name, found := range mobileNames {
		if !found {
			t.Errorf("/celebrations/2026 should include mobile celebration %q", name)
		}
	}
}

// CORS tests

func setupCORSMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /calendar/{year}", handleCalendar)
	return corsMiddleware(mux)
}

func TestCORSHeadersOnGET(t *testing.T) {
	handler := setupCORSMux()

	req := httptest.NewRequest("GET", "/calendar/2026", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", origin, "*")
	}
}

func TestCORSPreflightOptions(t *testing.T) {
	handler := setupCORSMux()

	req := httptest.NewRequest("OPTIONS", "/calendar/2026", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", origin, "*")
	}

	methods := w.Header().Get("Access-Control-Allow-Methods")
	if !strings.Contains(methods, "GET") {
		t.Errorf("Access-Control-Allow-Methods = %q, should contain GET", methods)
	}
}

// Phase 2 API tests

func setupCalendarMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /calendar/{year}/mobile-dates", handleMobileDates)
	mux.HandleFunc("GET /calendar/{year}/{month}/{day}", handleCalendarDay)
	mux.HandleFunc("GET /calendar/{year}", handleCalendar)
	return mux
}

func TestHandleCalendar200(t *testing.T) {
	mux := setupCalendarMux()

	req := httptest.NewRequest("GET", "/calendar/2026", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var entries []kalendar.CalendarEntry
	if err := json.Unmarshal(w.Body.Bytes(), &entries); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(entries) != 365 {
		t.Errorf("expected 365 entries, got %d", len(entries))
	}
}

func TestHandleCalendar400(t *testing.T) {
	mux := setupCalendarMux()

	req := httptest.NewRequest("GET", "/calendar/abc", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleCalendarDay200(t *testing.T) {
	mux := setupCalendarMux()

	req := httptest.NewRequest("GET", "/calendar/2026/1/1", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var entry kalendar.CalendarEntry
	if err := json.Unmarshal(w.Body.Bytes(), &entry); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if entry.Date != kalendar.NewDate(1, kalendar.JANUARY, 2026) {
		t.Errorf("date = %v, want 2026-01-01", entry.Date)
	}

	if entry.Season == "" {
		t.Error("season should not be empty")
	}

	if len(entry.Celebrations) == 0 {
		t.Error("Jan 1 should have celebrations")
	}
}

func TestHandleCalendarDay400InvalidMonth(t *testing.T) {
	mux := setupCalendarMux()

	req := httptest.NewRequest("GET", "/calendar/2026/13/1", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleCalendarDay400InvalidDay(t *testing.T) {
	mux := setupCalendarMux()

	req := httptest.NewRequest("GET", "/calendar/2026/1/32", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleCalendarDay400InvalidDate(t *testing.T) {
	mux := setupCalendarMux()

	tests := []struct {
		name string
		path string
	}{
		{"feb 30", "/calendar/2026/2/30"},
		{"feb 29 non-leap", "/calendar/2025/2/29"},
		{"apr 31", "/calendar/2026/4/31"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestHandleYearRangeValidation(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /celebrations/{year}", handleCelebrations)
	mux.HandleFunc("GET /liturgical-year/{year}", handleLiturgicalYear)
	mux.HandleFunc("GET /calendar/{year}/mobile-dates", handleMobileDates)
	mux.HandleFunc("GET /calendar/{year}/{month}/{day}", handleCalendarDay)
	mux.HandleFunc("GET /calendar/{year}", handleCalendar)

	tests := []struct {
		name string
		path string
	}{
		{"celebrations year too low", "/celebrations/1962"},
		{"celebrations year too high", "/celebrations/10000"},
		{"calendar year too low", "/calendar/1900"},
		{"calendar day year too low", "/calendar/1900/1/1"},
		{"liturgical-year too low", "/liturgical-year/1900"},
		{"mobile-dates too low", "/calendar/1900/mobile-dates"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestHandleCalendarDayResponse(t *testing.T) {
	mux := setupCalendarMux()

	req := httptest.NewRequest("GET", "/calendar/2026/4/5", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if _, ok := resp["date"]; !ok {
		t.Error("response should contain 'date'")
	}
	if _, ok := resp["season"]; !ok {
		t.Error("response should contain 'season'")
	}
	if _, ok := resp["season_color"]; !ok {
		t.Error("response should contain 'season_color'")
	}
	if _, ok := resp["celebrations"]; !ok {
		t.Error("response should contain 'celebrations'")
	}
}

func TestHandleMobileDates200(t *testing.T) {
	mux := setupCalendarMux()

	req := httptest.NewRequest("GET", "/calendar/2026/mobile-dates", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var celebrations []kalendar.Celebration
	if err := json.Unmarshal(w.Body.Bytes(), &celebrations); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(celebrations) == 0 {
		t.Error("mobile dates should not be empty")
	}

	for _, c := range celebrations {
		if !c.IsMovable {
			t.Errorf("celebration %q should have IsMovable=true", c.Name)
		}
	}
}

func TestHandleMobileDates400(t *testing.T) {
	mux := setupCalendarMux()

	req := httptest.NewRequest("GET", "/calendar/abc/mobile-dates", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleCalendarEntriesHaveSeasons(t *testing.T) {
	mux := setupCalendarMux()

	req := httptest.NewRequest("GET", "/calendar/2026", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var entries []kalendar.CalendarEntry
	json.Unmarshal(w.Body.Bytes(), &entries)

	for _, e := range entries {
		if e.Season == "" {
			t.Errorf("date %v has empty season", e.Date)
			break
		}
		if e.SeasonColor == "" {
			t.Errorf("date %v has empty season color", e.Date)
			break
		}
	}
}
