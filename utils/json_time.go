package utils

import "time"

// JSONTime used for marshal my own time format
type JSONTime time.Time

// MarshalJSON used for marshal my own time format
func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(time.RFC1123) + `"`), nil
}
