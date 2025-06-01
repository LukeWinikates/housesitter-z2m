package queries

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/timeofday"
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
	On     bool
	Device *database.Device
	Time   timeofday.SecondsInDay
}

func (q *TransitionsInWindowQuery) Find(start, end time.Time) []Transition {
	schedules := q.store.All()

	startSecondsInDay := timeofday.TimeToSecondsInDay(start)
	endSecondsInDay := timeofday.TimeToSecondsInDay(end)

	var result []Transition

	for _, s := range schedules {
		if between(startSecondsInDay, endSecondsInDay, s.OnTime) {
			for _, d := range s.DeviceSettings {
				result = append(result, Transition{
					On:     true,
					Device: d.Device,
					Time:   s.OnTime,
				})
			}
		}
		if between(startSecondsInDay, endSecondsInDay, s.OffTime) {
			for _, d := range s.DeviceSettings {
				result = append(result, Transition{
					On:     false,
					Device: d.Device,
					Time:   s.OffTime,
				})
			}
		}
	}

	return result
}

func between(startSecondsInDay, endSecondsInDay, transitionTime timeofday.SecondsInDay) bool {
	return startSecondsInDay <= transitionTime && transitionTime <= endSecondsInDay
}
