package database

import (
	"gorm.io/gorm"
)

type dbStore struct {
	database *gorm.DB
}

func (store *dbStore) Find(id string) (*Schedule, error) {
	s := &Schedule{}
	result := store.database.Preload("DeviceSettings").Preload("DeviceSettings.Device").Find(s, "id = ?", id)
	return s, result.Error
}

func (store *dbStore) All() []*Schedule {
	var all []*Schedule
	store.database.Preload("DeviceSettings").Preload("DeviceSettings.Device").Find(&all)
	return all
}

func (store *dbStore) SaveChanges(id string, s *Schedule) error {
	return store.database.Save(s).Error
}

func NewDBStore(db *gorm.DB) (Store, error) {
	return &dbStore{
		database: db,
	}, nil
}
