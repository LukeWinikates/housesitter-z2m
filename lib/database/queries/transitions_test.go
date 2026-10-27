package queries

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/dayofweek"
	"LukeWinikates/january-twenty-five/lib/timeofday"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type fakeStore struct {
	schedules []*database.Schedule
}

func (f fakeStore) Find(_ string) (*database.Schedule, error) {
	panic("not used for this test")
}

func (f fakeStore) All() []*database.Schedule {
	return f.schedules
}

func (f fakeStore) SaveChanges(_ string, _ *database.Schedule) error {
	panic("not used for this test")
}

func TestTransitionsInWindowQuery_Find(t *testing.T) {
	fakeDevice := &database.Device{
		FriendlyName: "Fake Device",
		IEEEAddress:  "ieeeeee",
		ID:           "fakeDevice",
	}
	fmt.Println(dayofweek.Daily)
	store := &fakeStore{
		schedules: []*database.Schedule{
			{
				OnTime:     timeofday.NewSecondsInDay(9, 15, timeofday.PM),
				OffTime:    timeofday.NewSecondsInDay(10, 45, timeofday.PM),
				DaysOfWeek: dayofweek.Daily,
				DeviceSettings: []*database.DeviceSetting{
					{
						Device:     fakeDevice,
						Brightness: 100,
						Model:      gorm.Model{},
						ScheduleID: "one",
						DeviceID:   fakeDevice.ID,
					},
				},
				ID: "one",
			},
			{
				OnTime:     timeofday.NewSecondsInDay(10, 45, timeofday.PM),
				OffTime:    timeofday.NewSecondsInDay(10, 45, timeofday.NEXT_DAY),
				DaysOfWeek: dayofweek.Daily,
				ID:         "two",
				DeviceSettings: []*database.DeviceSetting{
					{
						Device:     fakeDevice,
						Brightness: 100,
						Model:      gorm.Model{},
						ScheduleID: "two",
						DeviceID:   fakeDevice.ID,
					},
				},
			},
		},
	}
	startTime := timeofday.NewSecondsInDay(10, 40, timeofday.PM).OnDay(time.Local, time.Now())
	endTime := startTime.Add(15 * time.Minute)
	transitions := NewTransitionsInWindowQuery(store).Find(startTime, endTime)
	assert.Len(t, transitions, 2)
}
