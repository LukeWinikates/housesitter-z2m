package api

import (
	"LukeWinikates/january-twenty-five/lib/schedule"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ScheduleDevicePOSTHandler(scheduleStore schedule.Store, deviceStore zigbee2mqtt.Store) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Default().Printf("handling %s\n", request.RequestURI)
		scheduleIdFromPath := request.PathValue("schedule_id")
		decoder := json.NewDecoder(request.Body)
		var requestBody DeviceSettingPUTRequestBody
		err := decoder.Decode(&requestBody)
		if err != nil {
			fmt.Println(err)
			writer.WriteHeader(500) // actually a malformed request error
			return
			// general error handler / logger
		}

		_, err = scheduleStore.Find(scheduleIdFromPath)

		if err != nil {
			writer.WriteHeader(404)
			return
		}

		d, err := deviceStore.Find(requestBody.Id)
		if err != nil {
			writer.WriteHeader(404)
			return
		}

		deviceSettings := &schedule.DeviceSetting{Device: d}
		err = requestBody.Apply(deviceSettings)
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		err = scheduleStore.AddDeviceSettings(scheduleIdFromPath, deviceSettings)
		if err != nil {
			writer.WriteHeader(500)
			return
		}

		writer.WriteHeader(204)
	}
}
