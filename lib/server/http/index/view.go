package index

import (
	"LukeWinikates/january-twenty-five/lib/schedule"
	"fmt"
	"html/template"
)

type GridSchedule struct {
	OnTime       schedule.SecondsInDay
	OffTime      schedule.SecondsInDay
	FriendlyName string
	Devices      []GridDeviceSettings
	Row          int
	ID           string
}

type GridDeviceSettings struct {
	RowNumber      int
	DisplayClasses string
	FriendlyName   string
	ID             string
	Brightness     uint8
	Color          string
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

type Legend struct {
	DisplayClasses string
	Style          template.HTMLAttr
	Title          string
}

type ViewGrid struct {
	Schedules   []GridSchedule
	Legends     []Legend
	GridClasses string
}

func Grid(list []*schedule.Schedule) ViewGrid {

	var legends = make([]Legend, 48)

	for i := 0; i < 48; i++ {
		title := ""
		if i%2 == 0 {
			hour := i / 2 % 12
			if hour == 0 {
				hour = 12
			}
			title = fmt.Sprintf("%d", hour)
		}
		legends[i] = Legend{
			DisplayClasses: "legend",
			Style:          template.HTMLAttr(fmt.Sprintf("style=\"grid-column-start:tick %d\"", i+1)),
			Title:          title,
		}
	}
	return ViewGrid{
		Schedules:   displaySchedules(list),
		Legends:     legends,
		GridClasses: "",
	}
}

func displayDevices(devices []*schedule.DeviceSetting) []GridDeviceSettings {
	gridDevices := make([]GridDeviceSettings, len(devices))
	for i, device := range devices {
		gridDevices[i] = GridDeviceSettings{
			RowNumber:    i + 1,
			FriendlyName: device.Device.FriendlyName,
			ID:           device.Device.ID,
			Brightness:   device.Brightness,
			Color:        device.Color,
		}
	}
	return gridDevices
}

func displaySchedules(schedules []*schedule.Schedule) []GridSchedule {
	var result []GridSchedule
	for i, s := range schedules {
		result = append(result, GridSchedule{
			ID:           s.ID,
			OnTime:       s.OnTime,
			OffTime:      s.OffTime,
			FriendlyName: s.FriendlyName,
			Row:          i + 1,
			Devices:      displayDevices(s.DeviceSettings),
		})
	}
	return result
}
