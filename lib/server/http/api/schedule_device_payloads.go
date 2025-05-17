package api

import "LukeWinikates/january-twenty-five/lib/schedule"

type DeviceSettingPUTRequestBody struct {
	Id string `json:"id"`
	DeviceSettingPOSTRequestBody
}

type DeviceSettingPOSTRequestBody struct {
	Brightness int    `json:"brightness"`
	Color      string `json:"color"`
}

func (body DeviceSettingPOSTRequestBody) Apply(s *schedule.DeviceSetting) error {
	s.Brightness = uint8(body.Brightness)
	s.Color = body.Color
	return nil
}
