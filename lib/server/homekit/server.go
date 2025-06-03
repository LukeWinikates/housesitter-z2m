package homekit

import (
	"context"
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Start()
}

type s struct {
	hap *hap.Server
}

func (s s) Start() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		// Stop delivering signals.
		signal.Stop(c)
		// Cancel the context to stop the server.
		cancel()
	}()
	go func() {
		log.Default().Println(s.hap.ListenAndServe(ctx).Error())
	}()
}

func NewServer() (Server, error) {
	fs := hap.NewFsStore("tmp/hapdb")
	onOffSwitch := accessory.NewSwitch(accessory.Info{Name: "Rusuban Mode"})
	onOffSwitch.Switch.On.OnValueRemoteUpdate(func(on bool) {
		if on {
			log.Println("Switch is on")
		} else {
			log.Println("Switch is off")
		}
	})
	server, err := hap.NewServer(fs, onOffSwitch.A)
	server.Pin = "10010000"

	return &s{
		hap: server,
	}, err
}
