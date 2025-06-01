package server

import (
	"LukeWinikates/january-twenty-five/lib/schedule"
	"LukeWinikates/january-twenty-five/lib/server/http"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
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
}

func (r *realServer) Start() error {
	deviceChan, errChan := r.ztmClient.DeviceUpdates()
	go func() {
		for {
			select {
			case device := <-deviceChan:
				if device.Definition != nil && strings.Contains(device.Definition.Description, "bulb") {
					fmt.Println(device.FriendlyName)
					var savedDevice *schedule.Device
					r.database.Find(savedDevice, "ieee_address = ?", device.IEEEAddress)
					if savedDevice == nil {
						fmt.Printf("found new device: %s\n", device.FriendlyName)
						r.database.Create(&schedule.Device{
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

	fmt.Println("got here")
	return r.httpServer.Serve(r.options.Hostname)
}

func (r *realServer) Stop() error {
	return nil
}

type Options struct {
	DataDir  string
	Hostname string
}

func New(db *gorm.DB, client zigbee2mqtt.Client, opts Options) (Server, error) {
	return &realServer{
		database:   db,
		ztmClient:  client,
		options:    opts,
		httpServer: http.NewServer(db),
	}, nil
}
