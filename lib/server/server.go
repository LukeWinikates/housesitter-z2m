package server

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/runtime"
	"LukeWinikates/january-twenty-five/lib/server/homekit"
	"LukeWinikates/january-twenty-five/lib/server/http"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Server interface {
	Start() error
	Stop() error
}

type realServer struct {
	ztmClient  zigbee2mqtt.Client
	httpServer http.Server
	options    Options
	database   *gorm.DB
	hapServer  homekit.Server
}

func (r *realServer) Start() error {
	store, err := database.NewDBStore(r.database, false)
	if err != nil {
		return err
	}
	deviceChan, errChan := r.ztmClient.DeviceUpdates()
	go func() {
		fmt.Println("got here")
		for {
			fmt.Println("waiting for device message")
			select {
			case device := <-deviceChan:
				fmt.Println("device message")
				if device.Definition != nil && strings.Contains(device.Definition.Description, "bulb") {
					fmt.Println(device.FriendlyName)
					if r.database.Find(&database.Device{}, "ieee_address = ?", device.IEEEAddress).RowsAffected == 0 {
						fmt.Printf("found new device: %s\n", device.FriendlyName)
						r.database.Create(&database.Device{
							FriendlyName: device.FriendlyName,
							IEEEAddress:  device.IEEEAddress,
							ID:           uuid.New().String(),
						})
					}
				}
			case err := <-errChan:
				fmt.Println("received", err.Error())
			}
		}
	}()

	if r.options.Homekit {
		r.hapServer.Start()
	}

	runtime.NewRunner(store, r.ztmClient, r.options.Location).Start()
	return r.httpServer.Serve(r.options.Hostname)
}

func (r *realServer) Stop() error {
	return nil
}

type Options struct {
	DataDir  string
	Hostname string
	Homekit  bool
	Location *time.Location
}

func New(db *gorm.DB, client zigbee2mqtt.Client, opts Options) (Server, error) {
	hapServer, err := homekit.NewServer()
	if err != nil {
		return nil, err
	}
	return &realServer{
		database:   db,
		ztmClient:  client,
		options:    opts,
		httpServer: http.NewServer(db),
		hapServer:  hapServer,
	}, nil
}
