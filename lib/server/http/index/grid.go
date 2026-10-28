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
	legends = append(legends, Legend{
		DisplayClasses: "",
		Style:          template.HTMLAttr("style=\"grid-row-start: 1;grid-row-end: 1;grid-column-start:dayofweek 1;grid-column-end:dayofweek 7;\""),
		Title:          "Days of Week",
	})
	legends = append(legends, Legend{
		DisplayClasses: "",
		Style:          template.HTMLAttr("style=\"padding-right: 8px;grid-row-start: 1;grid-row-end: 1;grid-column-start:labels;grid-column-end:labels;\""),
		Title:          "Devices",
	})
	allGridDevices := toDeviceList(allDevices)
	return ViewGrid{
		Schedules:       displaySchedules(list, allGridDevices),
		EmptySchedule:   createEmptySchedule(allGridDevices),
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
	EmptySchedule   GridSchedule
}
