package kalendar

import (
	"testing"
)

func TestParseMonth(t *testing.T) {
	tests := []struct {
		input     string
		wantDay   int
		wantMonth Month
	}{
		{"2 de janeiro", 2, JANUARY},
		{"25 de dezembro", 25, DECEMBER},
		{"1 de junho", 1, JUNE},
		{"14 de fevereiro", 14, FEBRUARY},
		{"31 de maio", 31, MAY},
		{"29 de fevereiro", 29, FEBRUARY},
		{"10 de agosto", 10, AUGUST},
		{"17 de março", 17, MARCH},
		{"5 de novembro", 5, NOVEMBER},
		{"23 de abril", 23, APRIL},
		{"7 de outubro", 7, OCTOBER},
		{"30 de setembro", 30, SEPTEMBER},
		{"3 de julho", 3, JULY},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			day, month := parseMonth(tt.input)
			if day != tt.wantDay || month != tt.wantMonth {
				t.Errorf("parseMonth(%q) = (%d, %v), want (%d, %v)", tt.input, day, month, tt.wantDay, tt.wantMonth)
			}
		})
	}
}

func TestParseMonthInvalid(t *testing.T) {
	tests := []struct {
		input     string
		wantDay   int
		wantMonth Month
	}{
		{"", 1, JANUARY},
		{"  ", 1, JANUARY},
		{"janeiro", 1, JANUARY},
		{"invalid format", 1, JANUARY},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			day, month := parseMonth(tt.input)
			if day != tt.wantDay || month != tt.wantMonth {
				t.Errorf("parseMonth(%q) = (%d, %v), want (%d, %v)", tt.input, day, month, tt.wantDay, tt.wantMonth)
			}
		})
	}
}

func TestParseGrade(t *testing.T) {
	tests := []struct {
		input     string
		wantGrade CelebrationGrade
	}{
		{"Solenidade", GradeSolenidade},
		{"Festa", GradeFesta},
		{"Memória", GradeMemoria},
		{"Memória facultativa", GradeMemoriaFacultativa},
		{"Comemoração", GradeComemoracao},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			grade := parseGrade(tt.input)
			if grade != tt.wantGrade {
				t.Errorf("parseGrade(%q) = %v, want %v", tt.input, grade, tt.wantGrade)
			}
		})
	}
}

func TestParseGradeCaseInsensitive(t *testing.T) {
	tests := []struct {
		input     string
		wantGrade CelebrationGrade
	}{
		{"solenidade", GradeSolenidade},
		{"FESTA", GradeFesta},
		{"MEMÓRIA", GradeMemoria},
		{"memória facultativa", GradeMemoriaFacultativa},
		{"comemoração", GradeComemoracao},
		{"Solenidade", GradeSolenidade},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			grade := parseGrade(tt.input)
			if grade != tt.wantGrade {
				t.Errorf("parseGrade(%q) = %v, want %v", tt.input, grade, tt.wantGrade)
			}
		})
	}
}

func TestParseGradeDefault(t *testing.T) {
	grade := parseGrade("unknown")
	if grade != GradeMemoriaFacultativa {
		t.Errorf("parseGrade(%q) = %v, want %v", "unknown", grade, GradeMemoriaFacultativa)
	}
}

func TestCelebrationGradeLevel(t *testing.T) {
	tests := []struct {
		grade     CelebrationGrade
		wantLevel CelebrationLevel
	}{
		{GradeSolenidade, LevelSolenidade},
		{GradeFesta, LevelFesta},
		{GradeMemoria, LevelMemoria},
		{GradeMemoriaFacultativa, LevelMemoriaFacultativa},
		{GradeComemoracao, LevelComemoracao},
	}

	for _, tt := range tests {
		t.Run(string(tt.grade), func(t *testing.T) {
			level := tt.grade.Level()
			if level != tt.wantLevel {
				t.Errorf("%v.Level() = %v, want %v", tt.grade, level, tt.wantLevel)
			}
		})
	}
}

func TestIsFeastOfTheLord(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"Natal do Senhor", true},
		{"Epifania do Senhor", true},
		{"Batismo do Senhor", true},
		{"Anunciação do Senhor", true},
		{"Transfiguração do Senhor", true},
		{"Exaltação da Santa Cruz", true},
		{"Ascensão do Senhor", true},
		{"Santíssimo Corpo e Sangue de Cristo", true},
		{"Conversão de São Paulo", false},
		{"Cátedra de São Pedro", false},
		{"Natividade de São João Batista", false},
		{"Santos Pedro e Paulo", false},
		{"São Bartolomeu", false},
		{"São Mateus", false},
		{"São Lucas", false},
		{"São Marcos", false},
		{"Santo André", false},
		{"Santo Antônio", false},
		{"Santa Inês", false},
		{"São José", false},
		{"Santa Maria Madalena", false},
		{"Santos Anjos da Guarda", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isFeastOfTheLord(tt.name)
			if result != tt.expected {
				t.Errorf("isFeastOfTheLord(%q) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestLoadSaints(t *testing.T) {
	saints, err := loadSaints()
	if err != nil {
		t.Fatalf("loadSaints() error = %v", err)
	}

	if len(saints) == 0 {
		t.Error("loadSaints() returned empty slice")
	}

	for _, s := range saints {
		if s.Name == "" {
			t.Error("saint with empty Name")
		}
		if s.Date == "" {
			t.Error("saint with empty Date")
		}
		if s.Grade == "" {
			t.Error("saint with empty Grade")
		}
		if s.Color == "" {
			t.Error("saint with empty Color")
		}
		if s.Level == 0 {
			t.Error("saint with zero Level")
		}
	}
}

func TestLoadSaintsCount(t *testing.T) {
	saints, err := loadSaints()
	if err != nil {
		t.Fatalf("loadSaints() error = %v", err)
	}

	if len(saints) < 200 {
		t.Errorf("expected at least 200 saints, got %d", len(saints))
	}
}

func TestGetAllSaints(t *testing.T) {
	saints, err := GetAllSaints()
	if err != nil {
		t.Fatalf("GetAllSaints() error = %v", err)
	}

	if len(saints) == 0 {
		t.Error("GetAllSaints() returned empty slice")
	}
}

func TestSaintsCache(t *testing.T) {
	saints1, err1 := loadSaints()
	saints2, err2 := loadSaints()

	if err1 != nil || err2 != nil {
		t.Error("loadSaints() should not error on cached call")
	}

	if &saints1[0] == &saints2[0] {
		t.Log("Cache working - same slice reference")
	}
}

func TestGetSaintsForYear(t *testing.T) {
	saints, err := GetSaintsForYear(2026)
	if err != nil {
		t.Fatalf("GetSaintsForYear(2026) error = %v", err)
	}

	if len(saints) == 0 {
		t.Error("GetSaintsForYear() returned empty slice")
	}

	for _, s := range saints {
		if s.Name == "" {
			t.Error("saint with empty Name")
		}
	}
}

func TestGetLiturgicYearWithCelebrations(t *testing.T) {
	ly, err := GetLiturgicYearWithCelebrations(2026)
	if err != nil {
		t.Fatalf("GetLiturgicYearWithCelebrations(2026) error = %v", err)
	}

	if ly.MobileDates.Easter.Date.Year() != 2026 {
		t.Errorf("Easter year = %d, want 2026", ly.MobileDates.Easter.Date.Year())
	}

	if len(ly.Celebrations) == 0 {
		t.Error("celebrations slice is empty")
	}
}

func TestGetLiturgicYearWithCelebrationsMobileDates(t *testing.T) {
	tests := []struct {
		year int
		want Date
	}{
		{2025, NewDate(20, APRIL, 2025)},
		{2026, NewDate(5, APRIL, 2026)},
		{2027, NewDate(28, MARCH, 2027)},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.year)), func(t *testing.T) {
			ly, err := GetLiturgicYearWithCelebrations(tt.year)
			if err != nil {
				t.Fatalf("error = %v", err)
			}
			if ly.MobileDates.Easter.Date != tt.want {
				t.Errorf("Easter = %v, want %v", ly.MobileDates.Easter.Date, tt.want)
			}
		})
	}
}

func TestSaintFields(t *testing.T) {
	s := Saint{
		Name:             "Natal do Senhor",
		Date:             "25 de dezembro",
		Grade:            GradeSolenidade,
		Level:            LevelSolenidade,
		Color:            "branco",
		IsFeastOfTheLord: true,
	}

	if s.Name != "Natal do Senhor" {
		t.Errorf("Name = %v, want Natal do Senhor", s.Name)
	}
	if s.Grade != GradeSolenidade {
		t.Errorf("Grade = %v, want Solenidade", s.Grade)
	}
	if s.Level != LevelSolenidade {
		t.Errorf("Level = %v, want 1", s.Level)
	}
	if s.Color != "branco" {
		t.Errorf("Color = %v, want branco", s.Color)
	}
	if !s.IsFeastOfTheLord {
		t.Error("IsFeastOfTheLord should be true")
	}
}

func TestCelebrationFields(t *testing.T) {
	date := NewDate(25, DECEMBER, 2026)
	s := Saint{
		Name:  "Natal do Senhor",
		Date:  "25 de dezembro",
		Grade: GradeSolenidade,
		Color: "branco",
	}

	c := Celebration{
		Date:  date,
		Saint: s,
	}

	if c.Date != date {
		t.Errorf("Date = %v, want %v", c.Date, date)
	}
	if c.Saint.Name != "Natal do Senhor" {
		t.Errorf("Saint.Name = %v, want Natal do Senhor", c.Saint.Name)
	}
}

func TestLiturgicSeasonsWithCelebrationsFields(t *testing.T) {
	ly := &LiturgicSeasonsWithCelebrations{
		MobileDates: MobileDates{
			Easter: Feast{Date: NewDate(5, APRIL, 2026), Color: White},
		},
		LiturgicSeasons: LiturgicSeasons{
			Advent: Season{
				DateRange: DateRange{
					Start: NewDate(30, NOVEMBER, 2025),
					End:   NewDate(24, DECEMBER, 2025),
				},
				Color: Purple,
			},
		},
		Celebrations: []Celebration{},
	}

	if ly.MobileDates.Easter.Date.Year() != 2026 {
		t.Errorf("MobileDates.Easter.Year = %d, want 2026", ly.MobileDates.Easter.Date.Year())
	}
	if ly.LiturgicSeasons.Advent.Start.Month() != NOVEMBER {
		t.Errorf("Advent.Start.Month = %v, want NOVEMBER", ly.LiturgicSeasons.Advent.Start.Month())
	}
}
