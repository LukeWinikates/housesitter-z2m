package timeofday

import (
	"fmt"
	"time"
)

const Second SecondsInDay = 1
const Minute = 60 * Second
const Hour = 60 * Minute
const PM = 12 * Hour
const NEXT_DAY = 24 * Hour

type SecondsInDay int

func NewSecondsInDay(hour, minute, ampm SecondsInDay) SecondsInDay {
	return hour*Hour + minute*Minute + ampm
}

func (s SecondsInDay) HumanReadable() string {
	return fmt.Sprintf("%2d:%02d %s%s", s.Hour(), s.Minute(), s.AMPM(), s.DaysLater())
}

func (s SecondsInDay) Hour() int {
	hour := int(s/Hour) % 12
	if hour == 0 {
		return 12
	}
	return hour
}

func (s SecondsInDay) Minute() int {
	return int(s % Hour / Minute)
}

func (s SecondsInDay) HTMLValue() string {
	return fmt.Sprintf("%02d:%02d", s/Hour, s.Minute())
}

func (s SecondsInDay) AMPM() string {
	if s%NEXT_DAY >= PM {
		return "pm"
	}
	return "am"
}

func (s SecondsInDay) DaysLater() string {
	if s >= NEXT_DAY {
		return " (next day)"
	}
	return ""
}

func ToFriendlyTime(seconds SecondsInDay) string {
	duration, err := time.ParseDuration(fmt.Sprintf("%vs", seconds))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%v", duration)
}

func TimeToSecondsInDay(t time.Time) SecondsInDay {
	return SecondsInDay(
		(t.Hour() * int(Hour)) +
			(t.Minute() * int(Minute)) +
			(t.Second()),
	)
}

func (s SecondsInDay) Today(location *time.Location) time.Time {
	return s.OnDay(location, time.Now())
}

func (s SecondsInDay) OnDay(location *time.Location, onDay time.Time) time.Time {
	y, m, d := onDay.Date()
	return time.Date(y, m, d, s.Hour(), s.Minute(), 0, 0, location)
}
