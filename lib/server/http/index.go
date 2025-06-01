package http

import (
	"LukeWinikates/january-twenty-five/lib/schedule"
	"LukeWinikates/january-twenty-five/lib/server/http/index"
	"fmt"
	"net/http"
)

func indexPage(store schedule.Store, deviceStore schedule.DeviceStore) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		scheduleList := store.All()
		allDevices := deviceStore.All()
		viewModel := index.Grid(scheduleList, allDevices)
		err := homepageTemplate.Execute(writer, viewModel)
		if err != nil {
			writer.WriteHeader(500)
			fmt.Println(err.Error())
			return
		}
	}
}
