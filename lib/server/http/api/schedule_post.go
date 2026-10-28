package api

import (
	"LukeWinikates/january-twenty-five/lib/database/queries"
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"
)

func SchedulePOSTHandler(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var requestBody SchedulePOSTRequestBody
		err := decoder.Decode(&requestBody)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(500) // actually a malformed request error
			return
			// general error handler / logger
		}

		template, err := requestBody.ToScheduleCreateTemplate()
		if err != nil {
			log.Println(err)
			writer.WriteHeader(500)
			return
		}
		err = queries.CreateSchedule(db, *template)

		if err != nil {
			log.Println(err)
			writer.WriteHeader(500)
			return
		}

		writer.WriteHeader(204)
	}
}
