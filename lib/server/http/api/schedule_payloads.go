package api

import (
	"LukeWinikates/january-twenty-five/lib/schedule"
	"fmt"
	"strconv"
	"strings"
)

type SchedulePUTRequestBody struct {
	Id string `json:"id"`
	SchedulePOSTRequestBody
}

type SchedulePOSTRequestBody struct {
	Name    string `json:"name"`
	OnTime  string `json:"ontime"`
	OffTime string `json:"offtime"`
}

func (body SchedulePOSTRequestBody) Apply(s *schedule.Schedule) error {
	s.FriendlyName = body.Name
	onTime, err := htmlTimeToSecondsInDay(body.OnTime)
	if err != nil {
		return err
	}
	s.OnTime = onTime
	offTime, err := htmlTimeToSecondsInDay(body.OffTime)
	if err != nil {
		return err
	}
	s.OffTime = offTime
	return nil
}

// eg 22:15
func htmlTimeToSecondsInDay(time string) (schedule.SecondsInDay, error) {
	parts := strings.Split(time, ":")
	fmt.Println(parts)
	hrs, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	mins, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	return schedule.SecondsInDay((hrs * 60 * 60) + (mins * 60)), nil
}
