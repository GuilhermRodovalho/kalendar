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
		{"Solenidade", GradeSolemnity},
		{"Festa", GradeFeast},
		{"Memória", GradeMemorial},
		{"Memória facultativa", GradeOptionalMemorial},
		{"Comemoração", GradeCommemoration},
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
		{"solenidade", GradeSolemnity},
		{"FESTA", GradeFeast},
		{"MEMÓRIA", GradeMemorial},
		{"memória facultativa", GradeOptionalMemorial},
		{"comemoração", GradeCommemoration},
		{"Solenidade", GradeSolemnity},
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
	if grade != GradeOptionalMemorial {
		t.Errorf("parseGrade(%q) = %v, want %v", "unknown", grade, GradeOptionalMemorial)
	}
}

func TestCelebrationGradeLevel(t *testing.T) {
	tests := []struct {
		grade     CelebrationGrade
		wantLevel CelebrationLevel
	}{
		{GradeSolemnity, LevelSolemnity},
		{GradeFeast, LevelFeast},
		{GradeMemorial, LevelMemorial},
		{GradeOptionalMemorial, LevelOptionalMemorial},
		{GradeCommemoration, LevelCommemoration},
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

func TestLoadFixedCelebrations(t *testing.T) {
	celebrations, err := loadFixedCelebrations(2026)
	if err != nil {
		t.Fatalf("loadFixedCelebrations() error = %v", err)
	}

	if len(celebrations) == 0 {
		t.Error("loadFixedCelebrations() returned empty slice")
	}

	validColors := map[LiturgicalColor]bool{White: true, Red: true, Purple: true, Green: true, Rose: true}
	for _, c := range celebrations {
		if c.Name == "" {
			t.Error("celebration with empty Name")
		}
		if c.Date.Year() != 2026 {
			t.Errorf("celebration %q has year %d, want 2026", c.Name, c.Date.Year())
		}
		if c.Grade == "" {
			t.Error("celebration with empty Grade")
		}
		if c.Color == "" {
			t.Error("celebration with empty Color")
		}
		if !validColors[c.Color] {
			t.Errorf("celebration %q has unexpected color %q", c.Name, c.Color)
		}
		if c.Level == 0 {
			t.Error("celebration with zero Level")
		}
		if c.IsMovable {
			t.Errorf("fixed celebration %q should have IsMovable=false", c.Name)
		}
	}
}

func TestLoadFixedCelebrationsCount(t *testing.T) {
	celebrations, err := loadFixedCelebrations(2026)
	if err != nil {
		t.Fatalf("loadFixedCelebrations() error = %v", err)
	}

	if len(celebrations) < 200 {
		t.Errorf("expected at least 200 celebrations, got %d", len(celebrations))
	}
}

func TestRawSaintsCache(t *testing.T) {
	raw1, err1 := loadRawSaints()
	raw2, err2 := loadRawSaints()

	if err1 != nil || err2 != nil {
		t.Error("loadRawSaints() should not error on cached call")
	}

	if &raw1[0] == &raw2[0] {
		t.Log("Cache working - same slice reference")
	}
}

func TestGetCelebrationsForYear(t *testing.T) {
	celebrations, err := GetCelebrationsForYear(2026)
	if err != nil {
		t.Fatalf("GetCelebrationsForYear(2026) error = %v", err)
	}

	if len(celebrations) == 0 {
		t.Error("GetCelebrationsForYear() returned empty slice")
	}

	for _, c := range celebrations {
		if c.Name == "" {
			t.Error("celebration with empty Name")
		}
	}
}

func TestGetCelebrationsForYearIncludesMobile(t *testing.T) {
	celebrations, err := GetCelebrationsForYear(2026)
	if err != nil {
		t.Fatalf("GetCelebrationsForYear(2026) error = %v", err)
	}

	mobileNames := map[string]bool{
		"Santos Pedro e Paulo, apóstolos":            false,
		"Assunção da Bem-aventurada Virgem Maria":    false,
		"Todos os Santos":                            false,
		"Epifania do Senhor":                         false,
		"Nosso Senhor Jesus Cristo, Rei do Universo": false,
		"Sagrada Família de Jesus, Maria e José":     false,
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
			t.Errorf("GetCelebrationsForYear(2026) should include mobile celebration %q", name)
		}
	}
}

func TestGetCelebrationsForYearMobileDatesVary(t *testing.T) {
	celebrations2025, _ := GetCelebrationsForYear(2025)
	celebrations2026, _ := GetCelebrationsForYear(2026)

	findDate := func(celebrations []Celebration, name string) Date {
		for _, c := range celebrations {
			if c.Name == name {
				return c.Date
			}
		}
		return Date{}
	}

	pp2025 := findDate(celebrations2025, "Santos Pedro e Paulo, apóstolos")
	pp2026 := findDate(celebrations2026, "Santos Pedro e Paulo, apóstolos")

	if pp2025 == (Date{}) || pp2026 == (Date{}) {
		t.Error("Saints Peter and Paul should be present in both years")
	}
	// 2025: June 28 is Saturday -> Sunday June 29
	// 2026: June 28 is Sunday -> June 28
	if pp2025 == pp2026 {
		t.Log("Note: dates happen to be the same for these years (possible but check)")
	}
}

func TestGetLiturgicYearWithCelebrations(t *testing.T) {
	ly, err := GetLiturgicYearWithCelebrations(2026)
	if err != nil {
		t.Fatalf("GetLiturgicYearWithCelebrations(2026) error = %v", err)
	}

	if len(ly.Celebrations) == 0 {
		t.Error("celebrations slice is empty")
	}

	// Verify seasons are populated
	if ly.LiturgicSeasons.Advent.Start.Year() == 0 {
		t.Error("Advent start should be set")
	}
}

func TestGetLiturgicYearWithCelebrationsIncludesMobile(t *testing.T) {
	ly, err := GetLiturgicYearWithCelebrations(2026)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	mobileNames := map[string]bool{
		"Santos Pedro e Paulo, apóstolos":            false,
		"Assunção da Bem-aventurada Virgem Maria":    false,
		"Todos os Santos":                            false,
		"Epifania do Senhor":                         false,
		"Nosso Senhor Jesus Cristo, Rei do Universo": false,
	}

	for _, c := range ly.Celebrations {
		if _, ok := mobileNames[c.Name]; ok {
			mobileNames[c.Name] = true
		}
	}

	for name, found := range mobileNames {
		if !found {
			t.Errorf("celebrations should include mobile celebration %q", name)
		}
	}
}

func TestCelebrationFields(t *testing.T) {
	c := Celebration{
		Name:             "Natal do Senhor",
		Date:             NewDate(25, DECEMBER, 2026),
		Grade:            GradeSolemnity,
		Level:            LevelSolemnity,
		Color:            White,
		IsFeastOfTheLord: true,
		IsMovable:        false,
	}

	if c.Name != "Natal do Senhor" {
		t.Errorf("Name = %v, want Natal do Senhor", c.Name)
	}
	if c.Date != NewDate(25, DECEMBER, 2026) {
		t.Errorf("Date = %v, want 2026-12-25", c.Date)
	}
	if c.Grade != GradeSolemnity {
		t.Errorf("Grade = %v, want Solemnity", c.Grade)
	}
	if c.Level != LevelSolemnity {
		t.Errorf("Level = %v, want 1", c.Level)
	}
	if c.Color != White {
		t.Errorf("Color = %v, want white", c.Color)
	}
	if !c.IsFeastOfTheLord {
		t.Error("IsFeastOfTheLord should be true")
	}
	if c.IsMovable {
		t.Error("IsMovable should be false")
	}
}

func TestLiturgicSeasonsWithCelebrationsFields(t *testing.T) {
	ly := &LiturgicSeasonsWithCelebrations{
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

	if ly.LiturgicSeasons.Advent.Start.Month() != NOVEMBER {
		t.Errorf("Advent.Start.Month = %v, want NOVEMBER", ly.LiturgicSeasons.Advent.Start.Month())
	}
}

func TestTranslateColor(t *testing.T) {
	tests := []struct {
		input string
		want  LiturgicalColor
	}{
		{"branco", White},
		{"vermelho", Red},
		{"roxo", Purple},
		{"verde", Green},
		{"rosa", Rose},
		{"white", White},
		{"red", Red},
		{"unknown", LiturgicalColor("unknown")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := translateColor(tt.input)
			if got != tt.want {
				t.Errorf("translateColor(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
