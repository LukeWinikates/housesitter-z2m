package zigbee2mqtt

import (
	"LukeWinikates/january-twenty-five/lib/schedule"
	"fmt"
)

type Store interface {
	Find(id string) (*schedule.Device, error)
	All() []*schedule.Device
}

type inMemoryStore struct {
	devices []*schedule.Device
}

func (i inMemoryStore) Find(id string) (*schedule.Device, error) {
	for _, device := range i.devices {
		if device.ID == id {
			return device, nil
		}
	}
	return nil, fmt.Errorf("no device found for id: %s", id)
}

func (i inMemoryStore) All() []*schedule.Device {
	return i.devices
}

func NewInMemoryStore() Store {
	return &inMemoryStore{devices: []*schedule.Device{
		{
			FriendlyName: "Bedroom Nightstand",
			ID:           "3265D1FD-4FE5-4662-8AFE-C966089BCCB9",
		},
		{
			FriendlyName: "Playroom Light",
			ID:           "3265D1FD-4FE5-4662-8AFE-C9660891CCB9",
		},
		{
			FriendlyName: "Work Desk Lamp",
			ID:           "3265D1FD-4FE5-4662-8AFE-C9D6089BCCA9",
		},
		{
			FriendlyName: "Family Desk Lamp",
			ID:           "3265D1FD-4FE5-7662-8AFE-C966089BCCB9",
		},
	}}
}
