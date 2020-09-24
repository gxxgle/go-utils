package time

import (
	"time"
)

const (
	ISO8601      = "2006-01-02T15:04:05-0700"
	ISO8601Milli = "2006-01-02T15:04:05.999-0700"
)

// time.Month
const (
	M1  = time.January
	M2  = time.February
	M3  = time.March
	M4  = time.April
	M5  = time.May
	M6  = time.June
	M7  = time.July
	M8  = time.August
	M9  = time.September
	M10 = time.October
	M11 = time.November
	M12 = time.December
)

type Duration = time.Duration

// time.Duration
const (
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	Hour        = time.Hour
	Day         = time.Hour * 24
)

var (
	Sleep = time.Sleep
	Local = time.Local
)

type Time time.Time

func New(t time.Time) Time {
	return Time(t)
}

func Now() Time {
	return New(time.Now())
}

func Date(year int, month time.Month, day int, hour int, min int, sec int, nsec int) Time {
	return New(time.Date(year, month, day, hour, min, sec, nsec, time.Local))
}

func (t Time) Time() time.Time {
	return time.Time(t)
}

func (t Time) Unix() int64 {
	return t.Time().Unix()
}

func (t Time) UnixMilli() int64 {
	return t.Time().UnixNano() / int64(Millisecond)
}

func (t Time) Add(d Duration) Time {
	return New(t.Time().Add(d))
}

func (t Time) String() string {
	return t.Time().Format(ISO8601Milli)
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(ISO8601Milli)+2)
	b = append(b, '"')
	b = t.Time().AppendFormat(b, ISO8601Milli)
	b = append(b, '"')
	return b, nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	tm, err := time.Parse(`"`+ISO8601Milli+`"`, string(data))
	if err != nil {
		return err
	}

	*t = New(tm)
	return nil
}

func SetUTC() error {
	Local = time.UTC
	time.Local = time.UTC
	return nil
}

func SetChinaLocal() error {
	local, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		return err
	}

	Local = local
	time.Local = local
	return nil
}

func Since(t Time) Duration {
	return time.Since(t.Time())
}
