package kalendar

import (
	_ "embed"
	"encoding/json"
	"fmt"
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
		return GradeSolemnity
	case "festa":
		return GradeFeast
	case "memória", "memoria":
		return GradeMemorial
	case "memória facultativa", "memoria facultativa":
		return GradeOptionalMemorial
	case "comemoração", "comemoracao":
		return GradeCommemoration
	default:
		return GradeOptionalMemorial
	}
}

func translateColor(cor string) string {
	switch strings.ToLower(cor) {
	case "branco":
		return "white"
	case "vermelho":
		return "red"
	case "roxo":
		return "purple"
	case "verde":
		return "green"
	case "rosa":
		return "rose"
	default:
		return cor
	}
}

var lordFeasts = []string{
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

func isFeastOfTheLord(nome string) bool {
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
			Color:            translateColor(s.Cor),
			IsFeastOfTheLord: isFeastOfTheLord(s.Nome),
		}
	}

	saintsCache = saints
	return saintsCache, nil
}

func GetAllSaints() ([]Saint, error) {
	return loadSaints()
}

var monthNamesPortuguese = map[Month]string{
	JANUARY:   "janeiro",
	FEBRUARY:  "fevereiro",
	MARCH:     "março",
	APRIL:     "abril",
	MAY:       "maio",
	JUNE:      "junho",
	JULY:      "julho",
	AUGUST:    "agosto",
	SEPTEMBER: "setembro",
	OCTOBER:   "outubro",
	NOVEMBER:  "novembro",
	DECEMBER:  "dezembro",
}

// formatDatePortuguese formats a Date as "D de mês" in Portuguese.
func formatDatePortuguese(d Date) string {
	return fmt.Sprintf("%d de %s", d.Day(), monthNamesPortuguese[d.Month()])
}

// mobileSaints returns the mobile celebrations as Saints with dates resolved for the given year.
func mobileSaints(ly *LiturgicYear) []Saint {
	md := ly.MobileDates
	return []Saint{
		{
			Name:             "Epifania do Senhor",
			Date:             formatDatePortuguese(md.Epiphany.Date),
			Grade:            GradeSolemnity,
			Level:            LevelSolemnity,
			Color:            string(md.Epiphany.Color),
			IsFeastOfTheLord: true,
		},
		{
			Name:             "Batismo do Senhor",
			Date:             formatDatePortuguese(md.BaptismOfTheLord.Date),
			Grade:            GradeFeast,
			Level:            LevelFeast,
			Color:            string(md.BaptismOfTheLord.Color),
			IsFeastOfTheLord: true,
		},
		{
			Name:             "Santos Pedro e Paulo, apóstolos",
			Date:             formatDatePortuguese(md.SaintsPeterAndPaul.Date),
			Grade:            GradeSolemnity,
			Level:            LevelSolemnity,
			Color:            string(md.SaintsPeterAndPaul.Color),
			IsFeastOfTheLord: false,
		},
		{
			Name:             "Assunção da Bem-aventurada Virgem Maria",
			Date:             formatDatePortuguese(md.AssumptionOfMary.Date),
			Grade:            GradeSolemnity,
			Level:            LevelSolemnity,
			Color:            string(md.AssumptionOfMary.Color),
			IsFeastOfTheLord: false,
		},
		{
			Name:             "Todos os Santos",
			Date:             formatDatePortuguese(md.AllSaints.Date),
			Grade:            GradeSolemnity,
			Level:            LevelSolemnity,
			Color:            string(md.AllSaints.Color),
			IsFeastOfTheLord: false,
		},
		{
			Name:             "Nosso Senhor Jesus Cristo, Rei do Universo",
			Date:             formatDatePortuguese(md.ChristTheKing.Date),
			Grade:            GradeSolemnity,
			Level:            LevelSolemnity,
			Color:            string(md.ChristTheKing.Color),
			IsFeastOfTheLord: true,
		},
		{
			Name:             "Sagrada Família de Jesus, Maria e José",
			Date:             formatDatePortuguese(md.HolyFamily.Date),
			Grade:            GradeFeast,
			Level:            LevelFeast,
			Color:            string(md.HolyFamily.Color),
			IsFeastOfTheLord: true,
		},
		{
			Name:             "Bem-aventurada Virgem Maria, Mãe da Igreja",
			Date:             formatDatePortuguese(md.MaryMotherOfTheChurch.Date),
			Grade:            GradeMemorial,
			Level:            LevelMemorial,
			Color:            string(md.MaryMotherOfTheChurch.Color),
			IsFeastOfTheLord: false,
		},
		{
			Name:             "Imaculado Coração da Bem-aventurada Virgem Maria",
			Date:             formatDatePortuguese(md.ImmaculateHeartOfMary.Date),
			Grade:            GradeMemorial,
			Level:            LevelMemorial,
			Color:            string(md.ImmaculateHeartOfMary.Color),
			IsFeastOfTheLord: false,
		},
	}
}

// mobileCelebrations converts MobileDates fields into Celebrations.
func mobileCelebrations(ly *LiturgicYear) []Celebration {
	saints := mobileSaints(ly)
	year := ly.MobileDates.Epiphany.Date.Year()
	celebrations := make([]Celebration, len(saints))
	for i, s := range saints {
		day, month := parseMonth(s.Date)
		celebrations[i] = Celebration{
			Date:  NewDate(day, month, year),
			Saint: s,
		}
	}
	return celebrations
}

func GetSaintsForYear(year int) ([]Saint, error) {
	saints, err := loadSaints()
	if err != nil {
		return nil, err
	}

	ly := LiturgicYearOf(year)
	mobile := mobileSaints(ly)

	result := make([]Saint, len(saints), len(saints)+len(mobile))
	copy(result, saints)
	result = append(result, mobile...)

	return result, nil
}

func GetLiturgicYearWithCelebrations(year int) (*LiturgicSeasonsWithCelebrations, error) {
	ly := LiturgicYearOf(year)
	saints, err := loadSaints()
	if err != nil {
		return nil, err
	}

	mobile := mobileCelebrations(ly)
	celebrations := make([]Celebration, 0, len(saints)+len(mobile))
	for _, s := range saints {
		day, month := parseMonth(s.Date)
		date := NewDate(day, month, year)
		celebrations = append(celebrations, Celebration{
			Date:  date,
			Saint: s,
		})
	}
	celebrations = append(celebrations, mobile...)

	return &LiturgicSeasonsWithCelebrations{
		MobileDates:     ly.MobileDates,
		LiturgicSeasons: ly.LiturgicSeasons,
		Celebrations:    celebrations,
	}, nil
}
