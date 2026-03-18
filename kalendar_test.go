package kalendar

import (
	"reflect"
	"testing"
	"time"
)

func TestGetEasterByGauss(t *testing.T) {
	tests := []struct {
		input  int
		output time.Time
	}{
		{2026, time.Date(2026, time.April, 5, 0, 0, 0, 0, time.UTC)},
		{2025, time.Date(2025, time.April, 20, 0, 0, 0, 0, time.UTC)},
		{2024, time.Date(2024, time.March, 31, 0, 0, 0, 0, time.UTC)},
		{2023, time.Date(2023, time.April, 9, 0, 0, 0, 0, time.UTC)},
	}

	for _, input := range tests {
		result := EasterByGauss(input.input, GREGORIAN)
		if !result.Equal(input.output) {
			t.Errorf("Expected %v, received %v", input.output, result)
		}
	}
}

func TestGetLent(t *testing.T) {
	tests := []struct {
		input    int
		expected []time.Time
	}{
		{
			2026,
			[]time.Time{time.Date(2026, 2, 18, 0, 0, 0, 0, time.UTC), time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC)},
		},
	}

	for _, test := range tests {
		start, end := GetLent(test.input)

		if !reflect.DeepEqual([]time.Time{start, end}, test.expected) {
			t.Errorf("Expected %v, received %v", test.expected, []time.Time{start, end})
		}
	}
}

func TestLiturgicYearFromEaster(t *testing.T) {
	// Páscoa de 2026: 5 de abril
	ly := GetLiturgicYearOf(2026)

	cases := []struct {
		name     string
		got      Date
		expected Date
	}{
		{"ash_wednesday", ly.ash_wednesday, Date{18, FEBRUARY, 2026}},
		{"palm_sunday", ly.palm_sunday, Date{29, MARCH, 2026}},
		{"easter", ly.easter, Date{5, APRIL, 2026}},
		{"ascension_of_the_lord", ly.ascension_of_the_lord, Date{14, MAY, 2026}},
		{"pentecost", ly.pentecost, Date{24, MAY, 2026}},
		{"holy_trinity", ly.holy_trinity, Date{31, MAY, 2026}},
		{"corpus_christi", ly.corpus_christi, Date{4, JUNE, 2026}},
		{"feast_of_the_sacred_heart", ly.feast_of_the_sacred_heart, Date{12, JUNE, 2026}},
	}

	for _, c := range cases {
		if c.got != c.expected {
			t.Errorf("%s: got %+v, esperado %+v", c.name, c.got, c.expected)
		}
	}
}
