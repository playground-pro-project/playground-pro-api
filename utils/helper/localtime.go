package helper

import (
	"time"
)

type LocalTime time.Time

func (lt LocalTime) MarshalJSON() ([]byte, error) {
	t := time.Time(lt)
	if t.IsZero() {
		return []byte(`""`), nil
	}

	return []byte(`"` + t.Format("Monday, 02 January 2006 15:04:00") + `"`), nil
}
