package entities

import "time"

type RFC3339Time time.Time

func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	return []byte((time.Time(t)).Format("\"" + time.RFC3339 + "\"")), nil
}
