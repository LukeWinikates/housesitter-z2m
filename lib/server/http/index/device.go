package index

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"html/template"
)

type GridDeviceSettings struct {
	RowNumber    int
	FriendlyName string
	ID           string
	Brightness   uint8
	Color        string
	InSchedule   bool
	ScheduleID   string
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
	return "hidden hidable"
}
func toGridDeviceSettings(scheduleID string, devices []*database.DeviceSetting) []GridDeviceSettings {
	gridDevices := make([]GridDeviceSettings, len(devices))
	for i, device := range devices {
		gridDevices[i] = GridDeviceSettings{
			RowNumber:    i + 1,
			FriendlyName: device.Device.FriendlyName,
			ID:           device.Device.ID,
			Brightness:   device.Brightness,
			Color:        device.Color,
			InSchedule:   true,
			ScheduleID:   scheduleID,
		}
	}
	return gridDevices
}

type GridDevice struct {
	FriendlyName string
	ID           string
}

func (gd GridDevice) CreateEmptyDeviceSettings(scheduleID string) GridDeviceSettings {
	return GridDeviceSettings{
		FriendlyName: gd.FriendlyName,
		ID:           gd.ID,
		Brightness:   100,
		Color:        "#ffffff",
		ScheduleID:   scheduleID,
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
