package data

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// create a custom type for the release date to format user Dt
type Dt string

// handle custom defined release date for user...
func (t *Dt) UnmarshalJSON(jsonValue []byte) error {
	var ti string
	err := json.Unmarshal(jsonValue, &ti)

	cleanedDate := strings.Split(ti, ":")[0]
	if len(cleanedDate) != 10 {
		return errors.New("Invalid date format provided")
	}
	tt, err := time.Parse("2006-02-07", cleanedDate)
	if err != nil {
		return err
	}
	*t = Dt(tt.Format("2006-02-07"))
	return nil
}
