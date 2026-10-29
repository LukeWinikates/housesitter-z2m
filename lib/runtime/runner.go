package runtime

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/database/queries"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt"
	"log"
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
	changes   chan bool
}

func (r *runner) Changes() chan bool {
	return r.changes
}

func (r *runner) Start() {
	log.Println("Starting Schedule Runner")
	if r.running {
		log.Println("Runner was already running")
		return
	}
	r.running = true
	r.changes <- true

	r.lastRun = time.Now()
	r.scheduleEventsForNextInterval(time.Now().Add(r.interval))

	r.ticker = time.NewTicker(r.interval)
	go func() {
		for r.running {
			tick := <-r.ticker.C
			log.Printf("ticking at: %s\n", tick.String())
			r.scheduleEventsForNextInterval(tick.Add(r.interval))
		}
	}()
	log.Println("Running")
}

func (r *runner) Stop() {
	log.Println("Stopping Schedule Runner")
	if !r.running {
		log.Println("Runner was already stopped")
		return
	}
	r.running = false
	r.ticker.Stop()
	r.changes <- false
	log.Println("Stopped")
}
func (r *runner) Running() bool {
	return r.running
}

func (r *runner) scheduleEventsForNextInterval(end time.Time) {
	transitions := queries.NewTransitionsInWindowQuery(r.store).Find(r.lastRun, end)
	for _, t := range transitions {
		timer := time.NewTimer(time.Until(t.Time))
		go func() {
			<-timer.C
			msg := t.Message

			log.Printf("sending %v message for %s at %s\n", t.Kind, t.Device.FriendlyName, t.Time.Format(time.DateTime))
			err := r.z2mClient.SetDeviceState(t.Device.IEEEAddress, msg)
			if err != nil {
				log.Println("err")
			}
		}()
	}
	r.lastRun = time.Now()
}

type Runner interface {
	Start()
	Stop()
	Running() bool
	Changes() chan bool
}

func NewRunner(store database.Store, z2mClient zigbee2mqtt.Client, location *time.Location) Runner {
	return &runner{
		interval:  30 * time.Second,
		store:     store,
		z2mClient: z2mClient,
		location:  location,
		changes:   make(chan bool),
	}
}
