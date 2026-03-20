package kalendar

import (
	"encoding/json"
	"fmt"
)

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	var day, month, year int
	if _, err := fmt.Sscanf(s, "%d-%d-%d", &year, &month, &day); err != nil {
		return fmt.Errorf("invalid date format %q: expected YYYY-MM-DD", s)
	}
	*d = NewDate(day, Month(month), year)
	return nil
}
