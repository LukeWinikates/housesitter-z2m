package api

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"encoding/json"
	"fmt"
	"net/http"
)

func SchedulePUTHandler(scheduleStore database.Store) func(writer http.ResponseWriter, request *http.Request) {
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

		//err = requestBody.Apply(s)
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
