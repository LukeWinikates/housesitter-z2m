package api

import (
	"LukeWinikates/january-twenty-five/lib/schedule"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type SchedulePUTRequestBody struct {
	Id         string `json:"id"`
	OnTime     string `json:"ontime"`
	OffTime    string `json:"offtime"`
	Brightness string `json:"brightness"`
	Color      string `json:"color"`
}

func (body SchedulePUTRequestBody) Apply(sched *schedule.Schedule) error {
	ontime, err := htmlTimeToSecondsInDay(body.OnTime)
	if err != nil {
		return err
	}
	sched.OnTime = ontime
	offtime, err := htmlTimeToSecondsInDay(body.OffTime)
	if err != nil {
		return err
	}
	sched.OffTime = offtime
	return nil
}

// schedule device setting put

// eg 22:15
func htmlTimeToSecondsInDay(time string) (schedule.SecondsInDay, error) {
	parts := strings.Split(time, ":")
	fmt.Println(parts)
	hrs, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	fmt.Println(hrs)
	mins, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	fmt.Println(mins)
	return schedule.SecondsInDay((hrs * 60 * 60) + (mins * 60)), nil
}

func SchedulePutHandler(scheduleStore schedule.Store) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromPath := request.PathValue("schedule_id")
		fmt.Println(idFromPath)
		decoder := json.NewDecoder(request.Body)
		var requestBody SchedulePUTRequestBody
		err := decoder.Decode(&requestBody)
		if err != nil {
			fmt.Println(err)
			writer.WriteHeader(500) // actually a malformed request error
			return
			// general error handler / logger
		}
		if idFromPath != requestBody.Id {
			fmt.Println(err)
			fmt.Println("consistency issue...")
			writer.WriteHeader(500)
			return
		}

		s, err := scheduleStore.Find(idFromPath)
		if err != nil {
			writer.WriteHeader(404)
			return
		}

		err = requestBody.Apply(s)
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		err = scheduleStore.SaveChanges(idFromPath, s)
		if err != nil {
			writer.WriteHeader(500)
			return
		}

		writer.WriteHeader(204)
	}
}
