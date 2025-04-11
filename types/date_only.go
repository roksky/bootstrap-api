package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// DateOnly is a custom date type with 'yyyy-mm-dd' format
type DateOnly struct {
	time.Time
}

// NewDateOnly creates a new DateOnly instance
func NewDateOnly(t time.Time) DateOnly {
	// format data to 'yyyy-mm-dd'
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return DateOnly{t}
}

const dateLayout = "2006-01-02"

func (cd DateOnly) MarshalJSON() ([]byte, error) {
	if cd.IsZero() {
		return []byte("null"), nil
	}
	formatted := cd.Format("\"" + dateLayout + "\"")
	return []byte(formatted), nil
}

func (cd *DateOnly) UnmarshalJSON(b []byte) error {
	if b == nil {
		return nil
	}
	parsed, err := time.Parse("\""+dateLayout+"\"", string(b))
	if err != nil {
		return err
	}
	cd.Time = parsed
	return nil
}

func (cd DateOnly) Value() (driver.Value, error) {
	return cd.Format(dateLayout), nil
}

func (cd *DateOnly) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		cd.Time = v
		return nil
	case string:
		parsed, err := time.Parse(dateLayout, v)
		if err != nil {
			return err
		}
		cd.Time = parsed
		return nil
	default:
		return fmt.Errorf("can not convert %v to CustomDate", value)
	}
}
