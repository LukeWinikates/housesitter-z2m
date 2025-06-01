package api

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"encoding/json"
	"fmt"
	"net/http"
)

func SchedulePOSTHandler(scheduleStore database.Store) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var requestBody SchedulePOSTRequestBody
		err := decoder.Decode(&requestBody)
		if err != nil {
			fmt.Println(err)
			writer.WriteHeader(500) // actually a malformed request error
			return
			// general error handler / logger
		}

		newSchedule := &database.Schedule{}
		err = requestBody.Apply(newSchedule)

		if err != nil {
			writer.WriteHeader(500)
			return
		}
		err = scheduleStore.Add(newSchedule)

		if err != nil {
			writer.WriteHeader(500)
			return
		}

		writer.WriteHeader(204)
	}
}
