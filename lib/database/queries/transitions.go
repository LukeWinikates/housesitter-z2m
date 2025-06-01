package queries

import (
	"LukeWinikates/january-twenty-five/lib/database"
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

func (q *TransitionsInWindowQuery) Find(start, end time.Time) {
	schedules := q.store.All()

	startSecondsInDay :=

	for _, sched := range schedules {

	}

}
