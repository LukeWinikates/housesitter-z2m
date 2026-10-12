package index

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"fmt"
	"html/template"
)

func Grid(list []*database.Schedule, allDevices []*database.Device, isRunnerRunning bool) ViewGrid {

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
		Schedules:       displaySchedules(list, allGridDevices),
		AllDevices:      allGridDevices,
		Legends:         legends,
		GridClasses:     "",
		IsRunnerRunning: isRunnerRunning,
	}
}

type Legend struct {
	DisplayClasses string
	Style          template.HTMLAttr
	Title          string
}

type ViewGrid struct {
	Schedules       []GridSchedule
	Legends         []Legend
	GridClasses     string
	AllDevices      []GridDevice
	IsRunnerRunning bool
}
