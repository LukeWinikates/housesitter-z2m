package index

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/timeofday"
	"fmt"
	"html/template"
	"time"
)

type GridSchedule struct {
	OnTime           timeofday.SecondsInDay
	OffTime          timeofday.SecondsInDay
	FriendlyName     string
	Devices          []GridDeviceSettings
	Row              int
	ID               string
	AvailableDevices []GridDeviceSettings
	DaysOfWeek       []DayOfWeekCheckboxStates
}

func (s GridSchedule) FormattedTime() string {
	return fmt.Sprintf("%s - %s", s.OnTime.HumanReadable(), s.OffTime.HumanReadable())
}

func (s GridSchedule) InlineStyles() template.HTMLAttr {
	onTime := s.OnTime   // time in seconds
	offTime := s.OffTime // time in seconds
	// to column means -> 86400 second, divided by grid size 48
	columnSize := 86400 / 48
	// number of seconds as a half-hour
	startColumn := 1 + (int(onTime) / columnSize)
	endColumn := 1 + (int(offTime) / columnSize)
	return template.HTMLAttr(fmt.Sprintf("style=\"grid-row-start: %v; grid-column-start:tick %v ; grid-column-end: tick %v\"", s.Row+1, startColumn, endColumn))
}

func displaySchedules(schedules []*database.Schedule, allDevices []GridDevice) []GridSchedule {
	var result []GridSchedule
	for i, s := range schedules {
		settings := toGridDeviceSettings(s.DeviceSettings)
		var availableDevices []GridDeviceSettings
		for _, d := range allDevices {
			for _, set := range settings {
				if set.ID == d.ID {
					break
				}
			}
			availableDevices = append(availableDevices, d.CreateEmptyDeviceSettings())
		}

		result = append(result, GridSchedule{
			ID:               s.ID,
			OnTime:           s.OnTime,
			OffTime:          s.OffTime,
			DaysOfWeek:       daysOfWeekToCheckboxStates(s),
			FriendlyName:     s.FriendlyName,
			Row:              i + 1,
			Devices:          settings,
			AvailableDevices: availableDevices,
		})
	}
	return result
}

var intToDayOfWeek = map[int]string{
	0: "Sunday",
	1: "Monday",
	2: "Tuesday",
	3: "Wednesday",
	4: "Thursday",
	5: "Friday",
	6: "Saturday",
}

type DayOfWeekCheckboxStates struct {
	Name    string
	Checked template.HTMLAttr
}

func daysOfWeekToCheckboxStates(s *database.Schedule) []DayOfWeekCheckboxStates {
	var result []DayOfWeekCheckboxStates
	for i := 0; i < 7; i++ {
		checked := s.ActiveOnDayOfWeek(time.Weekday(i))
		checkedAttr := template.HTMLAttr("")
		if checked {
			checkedAttr = "checked"
		}
		result = append(result, DayOfWeekCheckboxStates{
			Name:    intToDayOfWeek[i],
			Checked: checkedAttr,
		})

	}
	return result
}

func (g ViewGrid) CreateEmptyDaysOfWeek() []DayOfWeekCheckboxStates {
	return daysOfWeekToCheckboxStates(&database.Schedule{})
}
