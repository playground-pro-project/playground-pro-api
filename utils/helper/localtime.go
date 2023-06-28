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

	return []byte(`"` + t.Format("2006-01-02 15:04:05") + `"`), nil
}
