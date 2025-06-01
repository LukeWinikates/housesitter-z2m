package server

import (
	"LukeWinikates/january-twenty-five/lib/server/http"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt/payloads"
	"fmt"
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
}

func (r *realServer) Start() error {
	r.ztmClient.SubscribeDeviceCatalog(func(devices []payloads.MessagePayload) {
		for _, device := range devices {
			fmt.Println(device.FriendlyName)
		}
	})
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
