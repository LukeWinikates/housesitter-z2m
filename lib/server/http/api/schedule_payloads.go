package api

import (
	"LukeWinikates/january-twenty-five/lib/database/queries"
	"LukeWinikates/january-twenty-five/lib/timeofday"
	"fmt"
	"strconv"
	"strings"
)

type SchedulePUTRequestBody struct {
	Id string `json:"id"`
	SchedulePOSTRequestBody
}

type DeviceSettingPayload struct {
	ID         string `json:"id"`
	Brightness int    `json:"brightness"`
	Color      string `json:"color"`
}

type SchedulePOSTRequestBody struct {
	Name    string                 `json:"name"`
	OnTime  string                 `json:"ontime"`
	OffTime string                 `json:"offtime"`
	Devices []DeviceSettingPayload `json:"devices"`
}

func (body SchedulePOSTRequestBody) ToScheduleCreateTemplate() (*queries.ScheduleCreateTemplate, error) {
	result := &queries.ScheduleCreateTemplate{}
	if body.Name != "" {
		result.Name = body.Name
	}
	onTime, err := htmlTimeToSecondsInDay(body.OnTime)
	if err != nil {
		return nil, err
	}
	offTime, err := htmlTimeToSecondsInDay(body.OffTime)
	if err != nil {
		return nil, err
	}

	return &queries.ScheduleCreateTemplate{
		Name:    body.Name,
		OnTime:  onTime,
		OffTime: offTime,
		Devices: mapDevices(body.Devices),
	}, nil
}

func mapDevices(devices []DeviceSettingPayload) []*queries.DeviceSettingCreate {
	var result []*queries.DeviceSettingCreate
	for _, device := range devices {
		result = append(result, &queries.DeviceSettingCreate{
			ID:         device.ID,
			Brightness: device.Brightness,
			Color:      device.Color,
		})
	}
	return result
}

// eg 22:15
func htmlTimeToSecondsInDay(time string) (timeofday.SecondsInDay, error) {
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
	return timeofday.SecondsInDay((hrs * 60 * 60) + (mins * 60)), nil
}
