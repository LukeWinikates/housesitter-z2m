package queries

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/timeofday"
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

type TransitionKind int

const (
	TransitionOn = iota
	TransitionOff
)

type Transition struct {
	Device  *database.Device
	Time    timeofday.SecondsInDay
	Message devices.LightControl
	Kind    TransitionKind
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
					Time:    timeofday.TimeToSecondsInDay(nextOn),
					Kind:    TransitionOn,
				})
			}
		}
		if nextOff := s.NextOffAfter(start); nextOff.Before(end) {
			for _, d := range s.DeviceSettings {
				result = append(result, Transition{
					Message: devices.OffMessage(),
					Device:  d.Device,
					Time:    timeofday.TimeToSecondsInDay(nextOff),
					Kind:    TransitionOff,
				})
			}
		}
	}

	return result
}
