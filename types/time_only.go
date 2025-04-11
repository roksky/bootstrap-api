package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// TimeOnly is a custom date type with 'yyyy-mm-dd' format
type TimeOnly struct {
	time.Time
}

// Layout for parsing and formatting time-only strings
const timeLayout = "15:04"

// UnmarshalJSON implements the json.Unmarshaler interface for TimeOnly
func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	if data == nil {
		return nil
	}
	// Remove the quotes from the JSON string
	str := string(data)
	str = str[1 : len(str)-1]

	// Parse the time-only string
	parsedTime, err := time.Parse(timeLayout, str)
	if err != nil {
		return err
	}

	// Set the parsed time to the TimeOnly struct
	t.Time = parsedTime
	return nil
}

// MarshalJSON implements the json.Marshaler interface for TimeOnly
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	// If the time is zero, return null
	if t.IsZero() {
		return json.Marshal(nil)
	}
	// Format the time to the specified layout
	formattedTime := t.Format(timeLayout)
	return json.Marshal(formattedTime)
}

func (t TimeOnly) Value() (driver.Value, error) {
	return t.Format(timeLayout), nil
}

func (t *TimeOnly) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	case string:
		parsed, err := time.Parse(timeLayout, v)
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	default:
		return fmt.Errorf("can not convert %v to TimeOnly", value)
	}
}
