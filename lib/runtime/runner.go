package runtime

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/database/queries"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt"
	"fmt"
	"time"
)

type runner struct {
	interval  time.Duration
	lastRun   time.Time
	store     database.Store
	z2mClient zigbee2mqtt.Client
	location  *time.Location
	running   bool
	ticker    *time.Ticker
}

func (r *runner) Start() {
	fmt.Println("Starting Schedule Runner")
	if r.running {
		fmt.Println("Runner was already running")
		return
	}
	r.running = true
	r.lastRun = time.Now()
	r.scheduleEventsForNextInterval(time.Now().Add(r.interval))

	r.ticker = time.NewTicker(r.interval)
	go func() {
		for r.running {
			tick := <-r.ticker.C
			fmt.Printf("ticking at: %s\n", tick.String())
			r.scheduleEventsForNextInterval(tick.Add(r.interval))
		}
	}()
	fmt.Println("Running")
}

func (r *runner) Stop() {
	fmt.Println("Stopping Schedule Runner")
	if !r.running {
		fmt.Println("Runner was already stopped")
		return
	}
	r.running = false
	r.ticker.Stop()
	fmt.Println("Stopped")
}
func (r *runner) Running() bool {
	return r.running
}

func (r *runner) scheduleEventsForNextInterval(end time.Time) {
	transitions := queries.NewTransitionsInWindowQuery(r.store).Find(r.lastRun, end)
	fmt.Printf("found %v transitions\n", len(transitions))
	for _, t := range transitions {
		timer := time.NewTimer(time.Until(t.Time.Today(r.location)))
		go func() {
			<-timer.C
			msg := t.Message

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
	Stop()
	Running() bool
}

func NewRunner(store database.Store, z2mClient zigbee2mqtt.Client, location *time.Location) Runner {
	return &runner{
		interval:  30 * time.Second,
		store:     store,
		z2mClient: z2mClient,
		location:  location,
	}
}
