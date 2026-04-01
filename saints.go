package kalendar

import (
	_ "embed"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
)

type saintsData struct {
	Dia  string `json:"dia"`
	Nome string `json:"nome"`
	Grau string `json:"grau"`
	Cor  string `json:"cor"`
}

//go:embed santos_missale.json
var santosData []byte

var (
	rawCache    []saintsData
	rawCacheErr error
	rawCacheOnce sync.Once
)

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

func translateColor(cor string) LiturgicalColor {
	switch strings.ToLower(cor) {
	case "branco", "white":
		return White
	case "vermelho", "red":
		return Red
	case "roxo", "purple":
		return Purple
	case "verde", "green":
		return Green
	case "rosa", "rose":
		return Rose
	default:
		return LiturgicalColor(cor)
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

func loadRawSaints() ([]saintsData, error) {
	rawCacheOnce.Do(func() {
		var rawData []saintsData
		if err := json.Unmarshal(santosData, &rawData); err != nil {
			rawCacheErr = err
			return
		}
		rawCache = rawData
	})
	return rawCache, rawCacheErr
}

func loadFixedCelebrations(year int) ([]Celebration, error) {
	raw, err := loadRawSaints()
	if err != nil {
		return nil, err
	}

	celebrations := make([]Celebration, len(raw))
	for i, s := range raw {
		day, month := parseMonth(s.Dia)
		grade := parseGrade(s.Grau)
		celebrations[i] = Celebration{
			Name:             s.Nome,
			Date:             NewDate(day, month, year),
			Grade:            grade,
			Level:            grade.Level(),
			Color:            translateColor(s.Cor),
			IsFeastOfTheLord: isFeastOfTheLord(s.Nome),
			IsMovable:        false,
		}
	}

	return celebrations, nil
}

func loadMobileCelebrations(ly *LiturgicYear) []Celebration {
	md := ly.MobileDates
	return []Celebration{
		{
			Name:             "Epifania do Senhor",
			Date:             md.Epiphany.Date,
			Grade:            GradeSolemnity,
			Level:            LevelSolemnity,
			Color:            md.Epiphany.Color,
			IsFeastOfTheLord: true,
			IsMovable:        true,
		},
		{
			Name:             "Batismo do Senhor",
			Date:             md.BaptismOfTheLord.Date,
			Grade:            GradeFeast,
			Level:            LevelFeast,
			Color:            md.BaptismOfTheLord.Color,
			IsFeastOfTheLord: true,
			IsMovable:        true,
		},
		{
			Name:             "Santos Pedro e Paulo, apóstolos",
			Date:             md.SaintsPeterAndPaul.Date,
			Grade:            GradeSolemnity,
			Level:            LevelSolemnity,
			Color:            md.SaintsPeterAndPaul.Color,
			IsFeastOfTheLord: false,
			IsMovable:        true,
		},
		{
			Name:             "Assunção da Bem-aventurada Virgem Maria",
			Date:             md.AssumptionOfMary.Date,
			Grade:            GradeSolemnity,
			Level:            LevelSolemnity,
			Color:            md.AssumptionOfMary.Color,
			IsFeastOfTheLord: false,
			IsMovable:        true,
		},
		{
			Name:             "Todos os Santos",
			Date:             md.AllSaints.Date,
			Grade:            GradeSolemnity,
			Level:            LevelSolemnity,
			Color:            md.AllSaints.Color,
			IsFeastOfTheLord: false,
			IsMovable:        true,
		},
		{
			Name:             "Nosso Senhor Jesus Cristo, Rei do Universo",
			Date:             md.ChristTheKing.Date,
			Grade:            GradeSolemnity,
			Level:            LevelSolemnity,
			Color:            md.ChristTheKing.Color,
			IsFeastOfTheLord: true,
			IsMovable:        true,
		},
		{
			Name:             "Sagrada Família de Jesus, Maria e José",
			Date:             md.HolyFamily.Date,
			Grade:            GradeFeast,
			Level:            LevelFeast,
			Color:            md.HolyFamily.Color,
			IsFeastOfTheLord: true,
			IsMovable:        true,
		},
		{
			Name:             "Bem-aventurada Virgem Maria, Mãe da Igreja",
			Date:             md.MaryMotherOfTheChurch.Date,
			Grade:            GradeMemorial,
			Level:            LevelMemorial,
			Color:            md.MaryMotherOfTheChurch.Color,
			IsFeastOfTheLord: false,
			IsMovable:        true,
		},
		{
			Name:             "Imaculado Coração da Bem-aventurada Virgem Maria",
			Date:             md.ImmaculateHeartOfMary.Date,
			Grade:            GradeMemorial,
			Level:            LevelMemorial,
			Color:            md.ImmaculateHeartOfMary.Color,
			IsFeastOfTheLord: false,
			IsMovable:        true,
		},
	}
}

func GetCelebrationsForYear(year int) ([]Celebration, error) {
	fixed, err := loadFixedCelebrations(year)
	if err != nil {
		return nil, err
	}

	ly := LiturgicYearOf(year)
	mobile := loadMobileCelebrations(ly)

	return append(fixed, mobile...), nil
}

func GetLiturgicYearWithCelebrations(year int) (*LiturgicSeasonsWithCelebrations, error) {
	ly := LiturgicYearOf(year)
	fixed, err := loadFixedCelebrations(year)
	if err != nil {
		return nil, err
	}

	mobile := loadMobileCelebrations(ly)

	return &LiturgicSeasonsWithCelebrations{
		LiturgicSeasons: ly.LiturgicSeasons,
		Celebrations:    append(fixed, mobile...),
	}, nil
}
