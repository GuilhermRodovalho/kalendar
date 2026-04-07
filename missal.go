package kalendar

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
)

// Preface represents a liturgical preface reference in the Missal.
type Preface struct {
	Name     string `json:"name"`
	Numeral  string `json:"numeral,omitempty"`
	Page     int    `json:"page"`
	Subtitle string `json:"subtitle,omitempty"`
}

// CommonReference points to a Missa dos Comuns section.
type CommonReference struct {
	Name string `json:"name"`
	Page int    `json:"page"`
}

// MissalReference holds the missal page and preface(s) for a celebration or day.
type MissalReference struct {
	MissalPage int              `json:"missal_page,omitempty"`
	Prefaces   []Preface        `json:"prefaces,omitempty"`
	CommonRef  *CommonReference `json:"common_ref,omitempty"`
	Note       string           `json:"note,omitempty"`
}

// missalDayEntry represents a single day or weekday group in the proprio_do_tempo JSON.
type missalDayEntry struct {
	MissalPage int       `json:"missal_page"`
	Prefaces   []Preface `json:"prefaces,omitempty"`
	Note       string    `json:"note,omitempty"`
}

// missalWeekEntry represents a week in the proprio_do_tempo JSON.
type missalWeekEntry struct {
	Domingo      *missalDayEntry `json:"domingo,omitempty"`
	Segunda      *missalDayEntry `json:"segunda,omitempty"`
	Terca        *missalDayEntry `json:"terca,omitempty"`
	Quarta       *missalDayEntry `json:"quarta,omitempty"`
	Quinta       *missalDayEntry `json:"quinta,omitempty"`
	Sexta        *missalDayEntry `json:"sexta,omitempty"`
	Sabado       *missalDayEntry `json:"sabado,omitempty"`
	DiasDeSemana *missalDayEntry `json:"dias_de_semana,omitempty"`
}

// missalSantosEntry represents a saint entry in the proprio_dos_santos JSON.
type missalSantosEntry struct {
	Name       string           `json:"name"`
	MissalPage int              `json:"missal_page"`
	Prefaces   []Preface        `json:"prefaces,omitempty"`
	CommonRef  *CommonReference `json:"common_ref,omitempty"`
}

// missalMovelEntry represents a movable celebration in the celebracoes_moveis JSON.
type missalMovelEntry struct {
	MissalPage int       `json:"missal_page"`
	Prefaces   []Preface `json:"prefaces,omitempty"`
}

// missalData represents the top-level structure of missal_references.json.
type missalData struct {
	ProprioDeTempo   map[string]map[string]missalWeekEntry `json:"proprio_do_tempo"`
	ProprioDosSantos map[string][]missalSantosEntry         `json:"proprio_dos_santos"`
	CelebMoveis      map[string]missalMovelEntry           `json:"celebracoes_moveis"`
}

//go:embed missal_references.json
var missalRefData []byte

var (
	missalCache     *missalData
	missalCacheErr  error
	missalCacheOnce sync.Once
)

func loadMissalReferences() (*missalData, error) {
	missalCacheOnce.Do(func() {
		var data missalData
		if err := json.Unmarshal(missalRefData, &data); err != nil {
			missalCacheErr = err
			return
		}
		missalCache = &data
	})
	return missalCache, missalCacheErr
}

// weekdayKey returns the JSON key for a given Weekday.
func weekdayKey(w Weekday) string {
	switch w {
	case SUNDAY:
		return "domingo"
	case MONDAY:
		return "segunda"
	case TUESDAY:
		return "terca"
	case WEDNESDAY:
		return "quarta"
	case THURSDAY:
		return "quinta"
	case FRIDAY:
		return "sexta"
	case SATURDAY:
		return "sabado"
	default:
		return ""
	}
}

// seasonKey returns the JSON key for a given SeasonName.
func seasonKey(s SeasonName) string {
	switch s {
	case SeasonAdvent:
		return "advento"
	case SeasonChristmas:
		return "natal"
	case SeasonOrdinaryTimeI, SeasonOrdinaryTimeII:
		return "tempo_comum"
	case SeasonLent:
		return "quaresma"
	case SeasonEasterTriduum:
		return "triduo_pascal"
	case SeasonEasterSeason:
		return "tempo_pascal"
	default:
		return "tempo_comum"
	}
}

// santosDateKey returns "MM-DD" format for looking up in proprio_dos_santos.
func santosDateKey(d Date) string {
	return fmt.Sprintf("%02d-%02d", d.month, d.day)
}

// resolveSeasonMissal looks up the season-level missal page and prefaces for a date.
func resolveSeasonMissal(d Date, season SeasonName, ly *LiturgicYear) (int, []Preface) {
	missal, err := loadMissalReferences()
	if err != nil || missal == nil {
		return 0, nil
	}

	sk := seasonKey(season)
	seasonData, ok := missal.ProprioDeTempo[sk]
	if !ok {
		return 0, nil
	}

	week := resolveWeekNumber(d, season, ly)
	weekKey := fmt.Sprintf("semana_%d", week)

	weekData, ok := seasonData[weekKey]
	if !ok {
		return 0, nil
	}

	wk := weekdayKey(d.Weekday())

	// Try specific day first
	entry := getDayEntry(&weekData, wk)
	if entry != nil {
		return entry.MissalPage, entry.Prefaces
	}

	// Fallback to dias_de_semana for weekdays
	if d.Weekday() != SUNDAY && weekData.DiasDeSemana != nil {
		return weekData.DiasDeSemana.MissalPage, weekData.DiasDeSemana.Prefaces
	}

	return 0, nil
}

// getDayEntry returns the day entry for a specific weekday key.
func getDayEntry(w *missalWeekEntry, key string) *missalDayEntry {
	switch key {
	case "domingo":
		return w.Domingo
	case "segunda":
		return w.Segunda
	case "terca":
		return w.Terca
	case "quarta":
		return w.Quarta
	case "quinta":
		return w.Quinta
	case "sexta":
		return w.Sexta
	case "sabado":
		return w.Sabado
	default:
		return nil
	}
}

// resolveWeekNumber calculates the week number within a liturgical season.
func resolveWeekNumber(d Date, season SeasonName, ly *LiturgicYear) int {
	switch season {
	case SeasonAdvent:
		adventStart := ly.LiturgicSeasons.Advent.Start
		days := int(d.toTime().Sub(adventStart.toTime()).Hours() / 24)
		return (days / 7) + 1

	case SeasonChristmas:
		christmasDay := NewDate(25, DECEMBER, ly.LiturgicSeasons.Christmas.Start.year)
		days := int(d.toTime().Sub(christmasDay.toTime()).Hours() / 24)
		if days < 0 {
			days = 0
		}
		return (days / 7) + 1

	case SeasonLent:
		ashWed := ly.MobileDates.AshWednesday.Date
		days := int(d.toTime().Sub(ashWed.toTime()).Hours() / 24)
		// Ash Wednesday through Saturday before 1st Sunday = week 0 (cinzas)
		if days < 4 {
			return 0
		}
		// 1st Sunday of Lent onwards
		firstSunday := ashWed.Next(SUNDAY)
		daysFromFirstSunday := int(d.toTime().Sub(firstSunday.toTime()).Hours() / 24)
		return (daysFromFirstSunday / 7) + 1

	case SeasonEasterTriduum:
		return 1

	case SeasonEasterSeason:
		easter := ly.MobileDates.Easter.Date
		days := int(d.toTime().Sub(easter.toTime()).Hours() / 24)
		return (days / 7) + 1

	case SeasonOrdinaryTimeI:
		baptism := ly.MobileDates.BaptismOfTheLord.Date
		start := baptism.Plus(1)
		days := int(d.toTime().Sub(start.toTime()).Hours() / 24)
		// OT I starts at week 1 (after Baptism of the Lord)
		return (days / 7) + 1

	case SeasonOrdinaryTimeII:
		// OT II resumes after Pentecost. The week number continues from where OT I left off.
		// Calculate which week of OT I was last, then continue.
		pentecost := ly.MobileDates.Pentecost.Date
		otIIStart := pentecost.Plus(1)
		daysIntoOTII := int(d.toTime().Sub(otIIStart.toTime()).Hours() / 24)

		// Calculate how many weeks of OT were used before Lent
		baptism := ly.MobileDates.BaptismOfTheLord.Date
		ashWed := ly.MobileDates.AshWednesday.Date
		otIDays := int(ashWed.toTime().Sub(baptism.Plus(1).toTime()).Hours() / 24)
		otIWeeks := (otIDays / 7) + 1

		// OT II starts at otIWeeks + 1
		return otIWeeks + (daysIntoOTII / 7) + 1

	default:
		return 1
	}
}

// resolveCelebrationMissal looks up missal reference for a specific celebration.
func resolveCelebrationMissal(name string, d Date, isMovable bool) *MissalReference {
	missal, err := loadMissalReferences()
	if err != nil || missal == nil {
		return nil
	}

	// Try movable celebrations first
	if isMovable {
		key := mobileCelebrationKey(name)
		if entry, ok := missal.CelebMoveis[key]; ok {
			return &MissalReference{
				MissalPage: entry.MissalPage,
				Prefaces:   entry.Prefaces,
			}
		}
	}

	// Try fixed santos
	dateKey := santosDateKey(d)
	if entries, ok := missal.ProprioDosSantos[dateKey]; ok {
		for _, entry := range entries {
			if nameMatchesSantos(name, entry.Name) {
				return &MissalReference{
					MissalPage: entry.MissalPage,
					Prefaces:   entry.Prefaces,
					CommonRef:  entry.CommonRef,
				}
			}
		}
	}

	return nil
}

// nameMatchesSantos checks if celebration name matches the santos entry.
func nameMatchesSantos(celebName, santosName string) bool {
	// Simple prefix matching — the santos_missale.json names may be longer
	if len(celebName) <= len(santosName) {
		return celebName == santosName || santosName[:len(celebName)] == celebName
	}
	return celebName[:len(santosName)] == santosName
}

// ResolveSeasonMissalForTest is exported for testing purposes.
func ResolveSeasonMissalForTest(d Date, season SeasonName, ly *LiturgicYear) (int, []Preface) {
	return resolveSeasonMissal(d, season, ly)
}

// mobileCelebrationKey maps celebration names to JSON keys in celebracoes_moveis.
func mobileCelebrationKey(name string) string {
	switch name {
	case "Quarta-feira de Cinzas":
		return "quarta_feira_de_cinzas"
	case "Domingo de Ramos da Paixão do Senhor":
		return "domingo_de_ramos"
	case "Sexta-feira Santa da Paixão do Senhor":
		return "sexta_feira_santa"
	case "Páscoa do Senhor":
		return "pascoa"
	case "Pentecostes":
		return "pentecostes"
	case "Santíssima Trindade":
		return "santissima_trindade"
	case "Santíssimo Corpo e Sangue de Cristo":
		return "corpus_christi"
	case "Sagrado Coração de Jesus":
		return "sagrado_coracao"
	case "Ascensão do Senhor":
		return "ascensao"
	case "Epifania do Senhor":
		return "epifania"
	case "Batismo do Senhor":
		return "batismo_do_senhor"
	case "Santos Pedro e Paulo, apóstolos":
		return "santos_pedro_e_paulo"
	case "Assunção da Bem-aventurada Virgem Maria":
		return "assuncao_de_maria"
	case "Todos os Santos":
		return "todos_os_santos"
	case "Nosso Senhor Jesus Cristo, Rei do Universo":
		return "cristo_rei"
	case "Sagrada Família de Jesus, Maria e José":
		return "sagrada_familia"
	case "Bem-aventurada Virgem Maria, Mãe da Igreja":
		return "mae_da_igreja"
	case "Imaculado Coração da Bem-aventurada Virgem Maria":
		return "imaculado_coracao_de_maria"
	default:
		return ""
	}
}
