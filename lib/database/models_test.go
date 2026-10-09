package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDaysOfWeek(t *testing.T) {
	s := &Schedule{}
	s.SetBooleanDaysOfWeek(false, true, true, true, true, true, false)
	assert.Equal(t, 62, int(s.DaysOfWeek))
	assert.ElementsMatch(t, []int{1, 2, 3, 4, 5}, s.NumeralDaysOfWeek())
	s.SetBooleanDaysOfWeek(false, true, true, false, false, false, false)
}

func TestSchedule_ActiveOnDayOfWeek(t *testing.T) {
	s := &Schedule{}
	s.SetBooleanDaysOfWeek(false, true, true, false, false, false, false)

	assert.True(t, s.ActiveOnDayOfWeek(time.Monday))
	assert.True(t, s.ActiveOnDayOfWeek(time.Tuesday))
	assert.False(t, s.ActiveOnDayOfWeek(time.Wednesday))
}

func TestSchedule_nextActiveDay(t *testing.T) {
	s := &Schedule{}
	s.SetBooleanDaysOfWeek(false, true, true, false, false, false, false)

	aTuesday := time.Date(2026, 10, 6, 0, 0, 0, 0, time.Local)

	assert.Equal(t, time.Monday, s.nextActiveDay(aTuesday).Weekday())
}
