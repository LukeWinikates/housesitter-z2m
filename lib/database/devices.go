package database

import (
	"gorm.io/gorm"
)

type DeviceStore interface {
	Find(id string) (*Device, error)
	All() []*Device
}

type dbDeviceStore struct {
	database *gorm.DB
}

func (d *dbDeviceStore) Find(id string) (*Device, error) {
	device := &Device{}
	result := d.database.Find(device, id)
	return device, result.Error
}

func (d *dbDeviceStore) All() []*Device {
	var all []*Device
	d.database.Find(&all)
	return all
}

func NewDBDeviceStore(db *gorm.DB) DeviceStore {
	return &dbDeviceStore{database: db}
}
