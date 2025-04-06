package schedule

import (
	"fmt"
	"time"
)

type HexColor = string

type Schedule struct {
	OnTime         SecondsInDay
	OffTime        SecondsInDay
	DeviceSettings []*DeviceSetting
	FriendlyName   string
	ID             string
}

type DeviceSetting struct {
	Device     *Device
	Brightness uint8
	Color      HexColor
}

type Device struct {
	FriendlyName string
	ID           string
}

const Second SecondsInDay = 1
const Minute = 60 * Second
const Hour = 60 * Minute
const PM = 12 * Hour

type SecondsInDay int

func TimeOfDay(hour, minute, ampm SecondsInDay) SecondsInDay {
	return hour*Hour + minute*Minute + ampm
}

func (s SecondsInDay) HumanReadable() string {
	return fmt.Sprintf("%2d:%02d %s", s.Hour(), s.Minute(), s.AMPM())
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
	if s >= PM {
		return "pm"
	}
	return "am"
}

func ToFriendlyTime(seconds SecondsInDay) string {
	duration, err := time.ParseDuration(fmt.Sprintf("%vs", seconds))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%v", duration)
}
