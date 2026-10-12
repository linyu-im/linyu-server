package localtime

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

var Location *time.Location = time.FixedZone("CST", 8*3600)

type LocalTime time.Time

const layout = "2006-01-02 15:04:05"

var scanLayouts = []string{
	layout,
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05.000000",
	"2006-01-02 15:04:05.000000000",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02 15:04:05.999999999Z07:00",
	time.RFC3339,
	time.RFC3339Nano,
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tt := time.Time(t).In(Location)
	str := fmt.Sprintf("\"%s\"", tt.Format(layout))
	return []byte(str), nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	if str == "" {
		*t = LocalTime(time.Time{})
		return nil
	}
	tt, err := time.ParseInLocation(layout, str, Location)
	if err != nil {
		return err
	}
	*t = LocalTime(tt)
	return nil
}

func (t LocalTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}
	// 固定无时区字符串，避免 SQLite 读回带 +08:00 导致 Scan 失败；MySQL 同样可接收
	return time.Time(t).In(Location).Format(layout), nil
}

func (t *LocalTime) Scan(value interface{}) error {
	if value == nil {
		*t = LocalTime(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = LocalTime(v.In(Location))
		return nil
	case []byte:
		return t.scanString(string(v))
	case string:
		return t.scanString(v)
	default:
		return fmt.Errorf("cannot scan type %T into LocalTime", value)
	}
}

func (t *LocalTime) scanString(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		*t = LocalTime(time.Time{})
		return nil
	}
	for _, layoutItem := range scanLayouts {
		if tt, err := time.ParseInLocation(layoutItem, s, Location); err == nil {
			*t = LocalTime(tt.In(Location))
			return nil
		}
		if tt, err := time.Parse(layoutItem, s); err == nil {
			*t = LocalTime(tt.In(Location))
			return nil
		}
	}
	return fmt.Errorf("cannot parse LocalTime: %q", s)
}

func (t LocalTime) ToTime() time.Time {
	return time.Time(t)
}

func Now() LocalTime {
	return LocalTime(time.Now().In(Location))
}
