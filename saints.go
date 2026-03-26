package kalendar

import (
	_ "embed"
	"encoding/json"
	"strconv"
	"strings"
)

type saintsData struct {
	Dia  string `json:"dia"`
	Nome string `json:"nome"`
	Grau string `json:"grau"`
	Cor  string `json:"cor"`
}

//go:embed santos_missale.json
var santosData []byte

var saintsCache []Saint

func parseMonth(dayStr string) (int, Month) {
	parts := strings.Fields(dayStr)
	if len(parts) < 3 {
		return 1, JANUARY
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil || day < 1 || day > 31 {
		return 1, JANUARY
	}

	monthStr := strings.ToLower(parts[2])
	var month Month

	switch monthStr {
	case "janeiro":
		month = JANUARY
	case "fevereiro":
		month = FEBRUARY
	case "março", "marco":
		month = MARCH
	case "abril":
		month = APRIL
	case "maio":
		month = MAY
	case "junho":
		month = JUNE
	case "julho":
		month = JULY
	case "agosto":
		month = AUGUST
	case "setembro":
		month = SEPTEMBER
	case "outubro":
		month = OCTOBER
	case "novembro":
		month = NOVEMBER
	case "dezembro":
		month = DECEMBER
	default:
		month = JANUARY
	}

	return day, month
}

func parseGrade(grau string) CelebrationGrade {
	switch strings.ToLower(grau) {
	case "solenidade":
		return GradeSolenidade
	case "festa":
		return GradeFesta
	case "memória", "memoria":
		return GradeMemoria
	case "memória facultativa", "memoria facultativa":
		return GradeMemoriaFacultativa
	case "comemoração", "comemoracao":
		return GradeComemoracao
	default:
		return GradeMemoriaFacultativa
	}
}

func isFeastOfTheLord(nome string) bool {
	lordFeasts := []string{
		"Natal do Senhor",
		"Epifania",
		"Batismo do Senhor",
		"Apresentação do Senhor",
		"Anunciação do Senhor",
		"Transfiguração do Senhor",
		"Exaltação da Santa Cruz",
		"Santíssimo Corpo e Sangue de Cristo",
		"Sagrado Coração de Jesus",
		"Santíssima Trindade",
		"Ascensão do Senhor",
		"Pentecostes",
	}

	for _, feast := range lordFeasts {
		if strings.Contains(nome, feast) {
			return true
		}
	}
	return false
}

func loadSaints() ([]Saint, error) {
	if saintsCache != nil {
		return saintsCache, nil
	}

	var rawData []saintsData
	if err := json.Unmarshal(santosData, &rawData); err != nil {
		return nil, err
	}

	saints := make([]Saint, len(rawData))
	for i, s := range rawData {
		grade := parseGrade(s.Grau)
		saints[i] = Saint{
			Name:             s.Nome,
			Date:             s.Dia,
			Grade:            grade,
			Level:            grade.Level(),
			Color:            s.Cor,
			IsFeastOfTheLord: isFeastOfTheLord(s.Nome),
		}
	}

	saintsCache = saints
	return saintsCache, nil
}

func GetAllSaints() ([]Saint, error) {
	return loadSaints()
}

func GetSaintsForYear(year int) ([]Saint, error) {
	saints, err := loadSaints()
	if err != nil {
		return nil, err
	}

	result := make([]Saint, len(saints))
	for i, s := range saints {
		result[i] = Saint{
			Name:             s.Name,
			Date:             s.Date,
			Grade:            s.Grade,
			Level:            s.Level,
			Color:            s.Color,
			IsFeastOfTheLord: s.IsFeastOfTheLord,
		}
	}

	return result, nil
}

func GetLiturgicYearWithCelebrations(year int) (*LiturgicSeasonsWithCelebrations, error) {
	ly := LiturgicYearOf(year)
	saints, err := loadSaints()
	if err != nil {
		return nil, err
	}

	celebrations := make([]Celebration, 0, len(saints))
	for _, s := range saints {
		day, month := parseMonth(s.Date)
		date := NewDate(day, month, year)
		celebrations = append(celebrations, Celebration{
			Date:  date,
			Saint: s,
		})
	}

	return &LiturgicSeasonsWithCelebrations{
		MobileDates:     ly.MobileDates,
		LiturgicSeasons: ly.LiturgicSeasons,
		Celebrations:    celebrations,
	}, nil
}
