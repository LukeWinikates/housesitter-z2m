package api

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

func ScheduleDELETEHandler(scheduleStore database.Store, db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromPath := request.PathValue("schedule_id")
		fmt.Println(idFromPath)
		s, err := scheduleStore.Find(idFromPath)
		if err != nil {
			writer.WriteHeader(404)
			return
		}

		if err = db.Delete(&s).Error; err != nil {
			fmt.Printf("err: %s\n", err.Error())
			writer.WriteHeader(500)
			return
		}

		writer.WriteHeader(204)
	}
}
