package homekit

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
)

type Server interface {
	Start()
	Stop()
}

type s struct {
	hap    *hap.Server
	cancel context.CancelFunc
}

func (s s) Start() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
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

func (s s) Stop() {
	s.cancel()
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
