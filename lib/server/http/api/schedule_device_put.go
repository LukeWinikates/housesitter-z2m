package api

import (
	"LukeWinikates/january-twenty-five/lib/schedule"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ScheduleDevicePUTHandler(scheduleStore schedule.Store) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Default().Printf("handling %s\n", request.RequestURI)
		scheduleIdFromPath := request.PathValue("schedule_id")
		deviceIdFromPath := request.PathValue("device_id")
		decoder := json.NewDecoder(request.Body)
		var requestBody DeviceSettingPUTRequestBody
		err := decoder.Decode(&requestBody)
		if err != nil {
			fmt.Println(err)
			writer.WriteHeader(500) // actually a malformed request error
			return
			// general error handler / logger
		}

		deviceSettings, err := scheduleStore.FindScheduleDeviceSettings(scheduleIdFromPath, deviceIdFromPath)
		if err != nil {
			writer.WriteHeader(404)
			return
		}

		err = requestBody.Apply(deviceSettings)
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		err = scheduleStore.SaveDeviceSettingChanges(scheduleIdFromPath, deviceIdFromPath, deviceSettings)
		if err != nil {
			writer.WriteHeader(500)
			return
		}

		writer.WriteHeader(204)
	}
}
