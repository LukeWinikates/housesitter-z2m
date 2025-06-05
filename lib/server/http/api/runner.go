package api

import (
	"LukeWinikates/january-twenty-five/lib/runtime"
	"encoding/json"
	"fmt"
	"net/http"
)

type RunnerPOSTBody struct {
	Run bool `json:"run"`
}

func RunnerStatePOSTHandler(runner runtime.Runner) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody *RunnerPOSTBody

		err := json.NewDecoder(request.Body).Decode(&requestBody)
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			writer.WriteHeader(500)

		}

		if requestBody.Run {
			runner.Start()
		} else {
			runner.Stop()
		}

		writer.WriteHeader(204)
	}
}
