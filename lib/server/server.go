package server

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/runtime"
	"LukeWinikates/january-twenty-five/lib/server/homekit"
	"LukeWinikates/january-twenty-five/lib/server/http"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	runner     runtime.Runner
}

func (r *realServer) Start() error {
	deviceChan, errChan := r.ztmClient.DeviceUpdates()
	go func() {
		for {
			log.Println("waiting for device message")
			select {
			case device := <-deviceChan:
				log.Println("received device message")
				if device.Definition != nil && strings.Contains(device.Definition.Description, "bulb") {
					log.Println(device.FriendlyName)
					if r.database.Find(&database.Device{}, "ieee_address = ?", device.IEEEAddress).RowsAffected == 0 {
						log.Printf("found new device: %s\n", device.FriendlyName)
						r.database.Create(&database.Device{
							FriendlyName: device.FriendlyName,
							IEEEAddress:  device.IEEEAddress,
							ID:           uuid.New().String(),
						})
					}
				}
			case err := <-errChan:
				log.Println("received", err.Error())
			}
		}
	}()

	log.Println("starting hap server")
	if r.options.Homekit {
		r.hapServer.Start()
	}

	log.Printf("starting schedule runner")
	r.runner.Start()
	log.Printf("starting http server at %s\n", r.options.Hostname)
	return r.httpServer.Serve(r.options.Hostname)
}

func (r *realServer) Stop() error {
	r.runner.Stop()
	r.hapServer.Stop()
	return nil
}

type Options struct {
	DataDir  string
	Hostname string
	Homekit  bool
	Location *time.Location
}

func New(db *gorm.DB, client zigbee2mqtt.Client, opts Options) (Server, error) {
	store := database.NewDBStore(db)
	runner := runtime.NewRunner(store, client, opts.Location)
	hapServer, err := homekit.NewServer(runner)
	if err != nil {
		return nil, err
	}
	return &realServer{
		database:   db,
		ztmClient:  client,
		options:    opts,
		httpServer: http.NewServer(db, runner),
		hapServer:  hapServer,
		runner:     runner,
	}, nil
}
