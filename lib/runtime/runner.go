package runtime

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/database/queries"
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
				fmt.Printf("ticking at: %s\n", tick.String())
				r.scheduleEventsForNextInterval(tick.Add(r.interval))
			}
		}
	}()
}

func (r *runner) scheduleEventsForNextInterval(end time.Time) {
	transitions := queries.NewTransitionsInWindowQuery(r.store).Find(r.lastRun, end)
	fmt.Printf("found %v transitions\n", len(transitions))
	for _, t := range transitions {
		timer := time.NewTimer(t.Time.Today().Sub(time.Now()))
		go func() {
			<-timer.C
			msg := devices.OffMessage()
			if t.On {
				msg = devices.OnMessage()
			}

			fmt.Printf("sending message for %s\n", t.Device.FriendlyName)
			err := r.z2mClient.SetDeviceState(t.Device.IEEEAddress, msg)
			if err != nil {
				fmt.Println("err")
			}
		}()
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
