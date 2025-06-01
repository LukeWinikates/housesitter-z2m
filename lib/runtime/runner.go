package runtime

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt/devices"
	"fmt"
	"time"
)

type runner struct {
	interval  time.Duration
	lastRun   time.Time
	store     database.Store
	z2mClient zigbee2mqtt.Client
}

func (r *runner) Start() {
	r.lastRun = time.Now()
	r.scheduleEventsForNextInterval(time.Now().Add(r.interval))

	ticker := time.NewTicker(r.interval)
	go func() {
		for {
			select {
			case tick := <-ticker.C:
				fmt.Printf("ticking at: %s", tick.String())
				r.scheduleEventsForNextInterval(tick.Add(r.interval))
			}
		}
	}()
}

func (r *runner) scheduleEventsForNextInterval(timer time.Time) {
	for _, s := range r.store.All() {
		// is the ontime between now and next interval?
		timer := &time.NewTimer(executionTime)
		go func() {
			<-timer.C
			for i, d := range s.DeviceSettings {
				d.Device.FriendlyName
				r.z2mClient.SetDeviceState(d.Device.FriendlyName, devices.OnMessage())
			}
		}()

		// is the offtime between now and next interval?
		//if s.OnTime >
	}
	r.lastRun = time.Now()

}

type Runner interface {
	Start()
}

func NewRunner(store database.Store, z2mClient zigbee2mqtt.Client) Runner {
	return &runner{
		interval:  30 * time.Second,
		store:     store,
		z2mClient: z2mClient,
	}
}
