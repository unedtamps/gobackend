package utils

import (
	"database/sql/driver"
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

var entropy = ulid.Monotonic(
	rand.New(rand.NewSource(time.Now().UnixNano())),
	0,
)

// ULID wraps ulid.ULID with database text (CHAR(26)) support
type ULID struct {
	ulid.ULID
}

// NewULID creates a new ULID
func NewULID() ULID {
	return ULID{ULID: ulid.MustNew(ulid.Timestamp(time.Now()), entropy)}
}

// ParseULID parses a ULID from string
func ParseULID(s string) (ULID, error) {
	u, err := ulid.Parse(s)
	if err != nil {
		return ULID{}, err
	}
	return ULID{ULID: u}, nil
}

// Scan implements sql.Scanner - reads text from database
func (u *ULID) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		parsed, err := ulid.Parse(v)
		if err != nil {
			return err
		}
		*u = ULID{ULID: parsed}
		return nil
	case []byte:
		return u.Scan(string(v))
	default:
		return fmt.Errorf("cannot scan %T into ULID", src)
	}
}

// Value implements driver.Valuer - writes text to database
func (u ULID) Value() (driver.Value, error) {
	return u.String(), nil
}

// GenerateULID generates a new ULID (backward compat)
func GenerateULID() ULID {
	return NewULID()
}
