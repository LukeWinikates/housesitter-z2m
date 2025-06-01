package api

import (
	"LukeWinikates/january-twenty-five/lib/database"
)

type DeviceSettingPUTRequestBody struct {
	Id string `json:"id"`
	DeviceSettingPOSTRequestBody
}

type DeviceSettingPOSTRequestBody struct {
	Brightness int    `json:"brightness"`
	Color      string `json:"color"`
}

func (body DeviceSettingPOSTRequestBody) Apply(s *database.DeviceSetting) error {
	s.Brightness = uint8(body.Brightness)
	s.Color = body.Color
	return nil
}
