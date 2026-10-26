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
	Devices          []GridDeviceSettings
	Row              int
	ID               string
	AvailableDevices []GridDeviceSettings
	DaysOfWeek       []DayOfWeekCheckboxStates
}

func (s GridSchedule) FormattedTime() string {
	return fmt.Sprintf("%s - %s", s.OnTime.HumanReadable(), s.OffTime.HumanReadable())
}

func (s GridSchedule) FriendlyName() string {
	return s.FormattedTime()
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

func (s GridSchedule) DayOfWeekText() string {
	var result = ""
	for _, d := range s.DaysOfWeek {
		if d.IsChecked {
			result = fmt.Sprintf("%s%s", result, d.Name[0:1])
		} else {
			result = fmt.Sprintf("%s%s", result, "-")
		}
	}
	return result
}

func displaySchedules(schedules []*database.Schedule, allDevices []GridDevice) []GridSchedule {
	var result []GridSchedule
	for i, s := range schedules {
		settings := toGridDeviceSettings(s.ID, s.DeviceSettings)
		var availableDevices []GridDeviceSettings
		for _, d := range allDevices {
			shouldAdd := true
			for _, set := range settings {
				if set.ID == d.ID {
					shouldAdd = false
					break
				}
			}
			if shouldAdd {
				availableDevices = append(availableDevices, d.CreateEmptyDeviceSettings(s.ID))
			}
		}

		result = append(result, GridSchedule{
			ID:               s.ID,
			OnTime:           s.OnTime,
			OffTime:          s.OffTime,
			DaysOfWeek:       daysOfWeekToCheckboxStates(s),
			Row:              i + 1,
			Devices:          settings,
			AvailableDevices: availableDevices,
		})
	}
	return result
}

func createEmptySchedule(allDevices []GridDevice) GridSchedule {
	var availableDevices []GridDeviceSettings
	for _, d := range allDevices {

		availableDevices = append(availableDevices, d.CreateEmptyDeviceSettings("new"))
	}

	return GridSchedule{
		ID:               "new",
		OnTime:           timeofday.NewSecondsInDay(9, 0, 0),
		OffTime:          timeofday.NewSecondsInDay(2, 0, timeofday.PM),
		DaysOfWeek:       createEmptyDaysOfWeek(),
		Row:              0,
		Devices:          nil,
		AvailableDevices: availableDevices,
	}
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
	Name      string
	IsChecked bool
	Checked   template.HTMLAttr
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
			Name:      intToDayOfWeek[i],
			IsChecked: checked,
			Checked:   checkedAttr,
		})
	}
	return result
}

func createEmptyDaysOfWeek() []DayOfWeekCheckboxStates {
	s := &database.Schedule{}
	s.SetBooleanDaysOfWeek(true, true, true, true, true, true, true)
	return daysOfWeekToCheckboxStates(s)
}
