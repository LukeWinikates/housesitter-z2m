package homekit

import (
	"LukeWinikates/january-twenty-five/lib/runtime"
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
	hap         *hap.Server
	cancel      context.CancelFunc
	onOffSwitch *accessory.Switch
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
	s.onOffSwitch.Switch.On.SetValue(false)
	s.cancel()
}

func NewServer(runner runtime.Runner) (Server, error) {
	fs := hap.NewFsStore("tmp/hapdb")
	onOffSwitch := accessory.NewSwitch(accessory.Info{Name: "Frokat Housesitter"})
	onOffSwitch.Switch.On.SetValue(runner.Running())
	go func(changes chan bool) {
		for {
			state := <-changes
			onOffSwitch.Switch.On.SetValue(state)
		}
	}(runner.Changes())
	onOffSwitch.Switch.On.OnValueRemoteUpdate(func(on bool) {
		if on {
			runner.Start()
			log.Println("Switch is on")
		} else {
			runner.Stop()
			log.Println("Switch is off")
		}
	})
	server, err := hap.NewServer(fs, onOffSwitch.A)
	//server.Addr = ":8999"
	if err != nil {
		return nil, err
	}
	server.Pin = "10010000"

	return &s{
		hap:         server,
		onOffSwitch: onOffSwitch,
	}, err
}
