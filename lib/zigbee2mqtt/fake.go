package zigbee2mqtt

import (
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt/devices"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt/payloads"
)

type noOpClient struct {
}

func (n noOpClient) DeviceUpdates() (chan payloads.MessagePayload, chan error) {
	return make(chan payloads.MessagePayload), make(chan error)
}

func (n noOpClient) SetDeviceState(_ string, _ devices.LightControl) error {
	return nil
}

func NoOpClient() Client {
	return noOpClient{}
}
