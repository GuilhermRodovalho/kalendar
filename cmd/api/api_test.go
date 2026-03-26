package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GuilhermRodovalho/kalendar"
)

func TestHandleLiturgicYear200(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /calendario-liturgico/{year}", handleLiturgicYear)

	req := httptest.NewRequest("GET", "/calendario-liturgico/2026", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp kalendar.LiturgicSeasonsWithCelebrations
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if resp.MobileDates.Easter.Date.Year() != 2026 {
		t.Errorf("Easter year = %d, want 2026", resp.MobileDates.Easter.Date.Year())
	}

	if len(resp.Celebrations) == 0 {
		t.Error("celebrations should not be empty")
	}
}

func TestHandleLiturgicYear400(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /calendario-liturgico/{year}", handleLiturgicYear)

	req := httptest.NewRequest("GET", "/calendario-liturgico/abc", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleSaints200(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /santos", handleSaints)

	req := httptest.NewRequest("GET", "/santos", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var saints []kalendar.Saint
	if err := json.Unmarshal(w.Body.Bytes(), &saints); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(saints) == 0 {
		t.Error("saints should not be empty")
	}

	for _, s := range saints {
		if s.Name == "" {
			t.Error("saint name should not be empty")
		}
	}
}

func TestHandleSaintsByYear200(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /santos/{year}", handleSaintsByYear)

	req := httptest.NewRequest("GET", "/santos/2025", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var saints []kalendar.Saint
	if err := json.Unmarshal(w.Body.Bytes(), &saints); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(saints) == 0 {
		t.Error("saints should not be empty")
	}
}

func TestHandleSaintsByYear400(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /santos/{year}", handleSaintsByYear)

	req := httptest.NewRequest("GET", "/santos/invalid", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCalendarLiturgicoResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /calendario-liturgico/{year}", handleLiturgicYear)

	req := httptest.NewRequest("GET", "/calendario-liturgico/2026", nil)
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

	if _, ok := resp["datas_moveis"]; !ok {
		t.Error("response should contain 'datas_moveis'")
	}
	if _, ok := resp["tempos_liturgicos"]; !ok {
		t.Error("response should contain 'tempos_liturgicos'")
	}
	if _, ok := resp["celebrações"]; !ok {
		t.Error("response should contain 'celebrações'")
	}
}

func TestMobileDatesNotEmpty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /calendario-liturgico/{year}", handleLiturgicYear)

	req := httptest.NewRequest("GET", "/calendario-liturgico/2026", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	var resp kalendar.LiturgicSeasonsWithCelebrations
	json.Unmarshal(w.Body.Bytes(), &resp)

	easter := kalendar.NewDate(5, kalendar.APRIL, 2026)
	if resp.MobileDates.Easter.Date != easter {
		t.Errorf("Easter = %v, want %v", resp.MobileDates.Easter.Date, easter)
	}
	if resp.MobileDates.AshWednesday.Date.Day() == 0 {
		t.Error("AshWednesday should not be zero day")
	}
	if resp.MobileDates.Pentecost.Date.Day() == 0 {
		t.Error("Pentecost should not be zero day")
	}
}

func TestCelebrationsNotEmpty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /calendario-liturgico/{year}", handleLiturgicYear)

	req := httptest.NewRequest("GET", "/calendario-liturgico/2026", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	var resp kalendar.LiturgicSeasonsWithCelebrations
	json.Unmarshal(w.Body.Bytes(), &resp)

	if len(resp.Celebrations) < 100 {
		t.Errorf("expected at least 100 celebrations, got %d", len(resp.Celebrations))
	}
}

func TestHandleLiturgicYearLeapYear(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /calendario-liturgico/{year}", handleLiturgicYear)

	req := httptest.NewRequest("GET", "/calendario-liturgico/2024", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp kalendar.LiturgicSeasonsWithCelebrations
	json.Unmarshal(w.Body.Bytes(), &resp)

	easter2024 := kalendar.NewDate(31, kalendar.MARCH, 2024)
	if resp.MobileDates.Easter.Date != easter2024 {
		t.Errorf("Easter 2024 = %v, want %v", resp.MobileDates.Easter.Date, easter2024)
	}
}

func TestHandleLiturgicYearVariousYears(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /calendario-liturgico/{year}", handleLiturgicYear)

	tests := []struct {
		name      string
		year      int
		wantDay   int
		wantMonth kalendar.Month
	}{
		{"2025", 2025, 20, kalendar.APRIL},
		{"2026", 2026, 5, kalendar.APRIL},
		{"2027", 2027, 28, kalendar.MARCH},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/calendario-liturgico/"+tt.name, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			var resp kalendar.LiturgicSeasonsWithCelebrations
			json.Unmarshal(w.Body.Bytes(), &resp)

			if resp.MobileDates.Easter.Date.Day() != tt.wantDay || resp.MobileDates.Easter.Date.Month() != tt.wantMonth {
				t.Errorf("Easter %d = %v, want %d-%v", tt.year, resp.MobileDates.Easter.Date, tt.wantDay, tt.wantMonth)
			}
		})
	}
}

func TestSaintsGradeDistribution(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /santos", handleSaints)

	req := httptest.NewRequest("GET", "/santos", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	var saints []kalendar.Saint
	json.Unmarshal(w.Body.Bytes(), &saints)

	grades := make(map[kalendar.CelebrationGrade]int)
	for _, s := range saints {
		grades[s.Grade]++
	}

	if grades[kalendar.GradeSolenidade] == 0 {
		t.Error("should have at least one Solenidade")
	}
	if grades[kalendar.GradeFesta] == 0 {
		t.Error("should have at least one Festa")
	}
	if grades[kalendar.GradeMemoria] == 0 {
		t.Error("should have at least one Memória")
	}
}

func TestSaintsColorDistribution(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /santos", handleSaints)

	req := httptest.NewRequest("GET", "/santos", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	var saints []kalendar.Saint
	json.Unmarshal(w.Body.Bytes(), &saints)

	colors := make(map[string]int)
	for _, s := range saints {
		colors[s.Color]++
	}

	if colors["branco"] == 0 {
		t.Error("should have saints with branco color")
	}
	if colors["vermelho"] == 0 {
		t.Error("should have saints with vermelho color")
	}
}
