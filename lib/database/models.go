package database

import (
	"LukeWinikates/january-twenty-five/lib/dayofweek"
	"LukeWinikates/january-twenty-five/lib/timeofday"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type HexColor = string

type Schedule struct {
	gorm.Model
	OnTime         timeofday.SecondsInDay
	OffTime        timeofday.SecondsInDay
	DaysOfWeek     dayofweek.DaysOfWeek
	DeviceSettings []*DeviceSetting
	FriendlyName   string
	ID             string `gorm:"primaryKey"`
}

func (s *Schedule) NumeralDaysOfWeek() []int {
	result := []int{}
	for i := uint8(0); i < 7; i++ {
		if s.DaysOfWeek>>(i)%2 == 1 {
			result = append(result, int(i))
		}
	}
	return result
}

func (s *Schedule) SetBooleanDaysOfWeek(sunday, monday, tuesday, wednesday, thursday, friday, saturday bool) {
	s.DaysOfWeek = 0
	if sunday {
		s.DaysOfWeek = s.DaysOfWeek | dayofweek.Sunday
	}
	if monday {
		s.DaysOfWeek = s.DaysOfWeek | dayofweek.Monday
	}
	if tuesday {
		s.DaysOfWeek = s.DaysOfWeek | dayofweek.Tuesday
	}
	if wednesday {
		s.DaysOfWeek = s.DaysOfWeek | dayofweek.Wednesday
	}
	if thursday {
		s.DaysOfWeek = s.DaysOfWeek | dayofweek.Thursday
	}
	if friday {
		s.DaysOfWeek = s.DaysOfWeek | dayofweek.Friday
	}
	if saturday {
		s.DaysOfWeek = s.DaysOfWeek | dayofweek.Saturday
	}
}

func (s *Schedule) NextOnAfter(afterTime time.Time) time.Time {
	return s.nextEventAfter(func() timeofday.SecondsInDay {
		return s.OnTime
	}, afterTime)
}

func (s *Schedule) ActiveOnDayOfWeek(weekday time.Weekday) bool {
	return s.DaysOfWeek>>(weekday)%2 == 1
}

func (s *Schedule) NextOffAfter(afterTime time.Time) time.Time {
	return s.nextEventAfter(func() timeofday.SecondsInDay { return s.OffTime }, afterTime)
}

func (s *Schedule) nextEventAfter(timeOfDayFunc func() timeofday.SecondsInDay, afterTime time.Time) time.Time {
	tod := timeOfDayFunc()
	if tod > timeofday.TimeToSecondsInDay(afterTime) && s.ActiveOnDayOfWeek(afterTime.Weekday()) {
		day := tod.OnDay(afterTime.Location(), afterTime)
		return day
	}
	return tod.OnDay(afterTime.Location(), s.nextActiveDay(afterTime))
}

func (s *Schedule) nextActiveDay(afterTime time.Time) time.Time {
	for daysToAdd := 1; daysToAdd < 8; daysToAdd++ {
		daysOffset := (afterTime.Weekday() + time.Weekday(daysToAdd)) % 7
		if s.ActiveOnDayOfWeek(daysOffset) {
			return afterTime.AddDate(0, 0, daysToAdd)
		}
	}
	fmt.Println("empty or invalid days of week, treating schedule as daily")
	return afterTime.AddDate(0, 0, 1)
}

type DeviceSetting struct {
	Device     *Device
	Brightness uint8
	Color      HexColor
	gorm.Model
	ScheduleID string
	DeviceID   string
}

type Device struct {
	gorm.Model
	FriendlyName string
	IEEEAddress  string
	ID           string `gorm:"primaryKey"`
}

type Settings struct {
	gorm.Model
	Active bool
}
