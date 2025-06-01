package database

import (
	"LukeWinikates/january-twenty-five/lib/time"
	"gorm.io/gorm"
)

type HexColor = string

type Schedule struct {
	gorm.Model
	OnTime         timeofday.SecondsInDay
	OffTime        timeofday.SecondsInDay
	DeviceSettings []*DeviceSetting
	FriendlyName   string
	ID             string `gorm:"primaryKey"`
}

type DeviceSetting struct {
	Device     *Device
	Brightness uint8
	Color      HexColor
	gorm.Model
	ScheduleID string
	DeviceID   string
}

type Device struct {
	gorm.Model
	FriendlyName string
	IEEEAddress  string
	ID           string `gorm:"primaryKey"`
}

type Settings struct {
	gorm.Model
	Active bool
}
