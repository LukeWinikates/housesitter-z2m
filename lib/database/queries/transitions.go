package queries

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt/devices"
	"time"
)

type TransitionsInWindowQuery struct {
	store database.Store
}

func NewTransitionsInWindowQuery(store database.Store) *TransitionsInWindowQuery {
	return &TransitionsInWindowQuery{
		store: store,
	}
}

type Transition struct {
	Device  *database.Device
	Time    time.Time
	Message devices.LightControl
	Kind    string
}

func (q *TransitionsInWindowQuery) Find(start, end time.Time) []Transition {
	schedules := q.store.All()

	var result []Transition

	for _, s := range schedules {
		if nextOn := s.NextOnAfter(start); nextOn.Before(end) {
			for _, d := range s.DeviceSettings {
				result = append(result, Transition{
					Message: devices.OnMessage(int(d.Brightness)),
					Device:  d.Device,
					Time:    nextOn,
					Kind:    "On",
				})
			}
		}
		if nextOff := s.NextOffAfter(start); nextOff.Before(end) {
			for _, d := range s.DeviceSettings {
				result = append(result, Transition{
					Message: devices.OffMessage(),
					Device:  d.Device,
					Time:    nextOff,
					Kind:    "Off",
				})
			}
		}
	}

	return result
}
