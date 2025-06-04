package database

import (
	"gorm.io/gorm"
)

type DeviceStore interface {
	All() []*Device
}

type dbDeviceStore struct {
	database *gorm.DB
}

func (d *dbDeviceStore) All() []*Device {
	var all []*Device
	d.database.Find(&all)
	return all
}

func NewDBDeviceStore(db *gorm.DB) DeviceStore {
	return &dbDeviceStore{database: db}
}
