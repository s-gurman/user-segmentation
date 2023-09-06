package t

import (
	"fmt"
	"strings"
	"time"
)

const Layout = "2006-01-02 15:04:05"

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(data []byte) (err error) {
	s := strings.Trim(string(data), `"`)
	t.Time, err = time.ParseInLocation(Layout, s, time.Local)
	return err
}

func (t CustomTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Format(Layout))), nil
}

func (t CustomTime) String() string {
	return t.Format(Layout)
}
