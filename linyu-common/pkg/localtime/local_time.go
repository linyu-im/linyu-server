package localtime

import (
	"database/sql/driver"
	"fmt"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"strings"
	"time"
)

var Location *time.Location = time.FixedZone("CST", 8*3600)

func InitLocalTime() {
	loc, err := time.LoadLocation(config.C.Mysql.Timezone)
	if err != nil {
		panic("failed to init time zone: " + err.Error())
	}
	Location = loc
}

type LocalTime time.Time

const layout = "2006-01-02 15:04:05"

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
	return time.Time(t).In(Location), nil
}

func (t *LocalTime) Scan(value interface{}) error {
	if value == nil {
		*t = LocalTime(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = LocalTime(v.In(Location))
	case []byte:
		tt, err := time.ParseInLocation(layout, string(v), Location)
		if err != nil {
			return err
		}
		*t = LocalTime(tt)
	case string:
		tt, err := time.ParseInLocation(layout, v, Location)
		if err != nil {
			return err
		}
		*t = LocalTime(tt)
	default:
		return fmt.Errorf("cannot scan type %T into LocalTime", value)
	}
	return nil
}

func (t LocalTime) ToTime() time.Time {
	return time.Time(t)
}

func Now() LocalTime {
	return LocalTime(time.Now().In(Location))
}
