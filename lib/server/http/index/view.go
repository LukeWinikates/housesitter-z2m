package index

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/timeofday"
	"fmt"
	"html/template"
)

type GridSchedule struct {
	OnTime           timeofday.SecondsInDay
	OffTime          timeofday.SecondsInDay
	FriendlyName     string
	Devices          []GridDeviceSettings
	Row              int
	ID               string
	AvailableDevices []GridDeviceSettings
}

type GridDeviceSettings struct {
	RowNumber int
	//DisplayClasses string
	FriendlyName string
	ID           string
	Brightness   uint8
	Color        string
	InSchedule   bool
}

func (d GridDeviceSettings) Checked() template.HTMLAttr {
	if d.InSchedule {
		return "checked"
	}
	return ""
}
func (d GridDeviceSettings) DisplayClasses() template.HTMLAttr {
	if d.InSchedule {
		return ""
	}
	return "hidden"
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
	AllDevices  []GridDevice
}

func Grid(list []*database.Schedule, allDevices []*database.Device) ViewGrid {

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
	allGridDevices := toDeviceList(allDevices)
	return ViewGrid{
		Schedules:   displaySchedules(list, allGridDevices),
		AllDevices:  allGridDevices,
		Legends:     legends,
		GridClasses: "",
	}
}

func toGridDeviceSettings(devices []*database.DeviceSetting) []GridDeviceSettings {
	gridDevices := make([]GridDeviceSettings, len(devices))
	for i, device := range devices {
		gridDevices[i] = GridDeviceSettings{
			RowNumber:    i + 1,
			FriendlyName: device.Device.FriendlyName,
			ID:           device.Device.ID,
			Brightness:   device.Brightness,
			Color:        device.Color,
			InSchedule:   true,
		}
	}
	return gridDevices
}

type GridDevice struct {
	FriendlyName string
	ID           string
}

func (gd GridDevice) CreateEmptyDeviceSettings() GridDeviceSettings {
	return GridDeviceSettings{
		FriendlyName: gd.FriendlyName,
		ID:           gd.ID,
		Brightness:   100,
		Color:        "#ffffff",
	}
}

func toDeviceList(devices []*database.Device) []GridDevice {
	var result []GridDevice
	for _, d := range devices {
		result = append(result, GridDevice{
			FriendlyName: d.FriendlyName,
			ID:           d.ID,
		})
	}
	return result
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
			FriendlyName:     s.FriendlyName,
			Row:              i + 1,
			Devices:          settings,
			AvailableDevices: availableDevices,
		})
	}
	return result
}
