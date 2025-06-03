package database

import (
	"LukeWinikates/january-twenty-five/lib/timeofday"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type dbStore struct {
	database *gorm.DB
}

func (store *dbStore) Add(schedule *Schedule) error {
	schedule.ID = uuid.New().String()
	result := store.database.Create(schedule)
	return result.Error
}

func (store *dbStore) AddDeviceSettings(scheduleID string, settings *DeviceSetting) error {
	s := &Schedule{}
	result := store.database.Find(s, scheduleID)
	if result.Error != nil {
		return result.Error
	}
	return store.database.Model(s).Association("DeviceSettings").Append(settings)
}

func (store *dbStore) SaveDeviceSettingChanges(scheduleId string, deviceId string, settings *DeviceSetting) error {
	return store.database.Save(settings).Error
}

func (store *dbStore) FindScheduleDeviceSettings(scheduleId, deviceId string) (*DeviceSetting, error) {
	schedule, err := store.Find(scheduleId)
	if err != nil {
		return nil, err
	}
	for _, setting := range schedule.DeviceSettings {
		if setting.Device.ID == deviceId {
			return setting, nil
		}
	}
	return nil, fmt.Errorf("no device found for %s in schedule %s", deviceId, scheduleId)
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

func NewDBStore(db *gorm.DB, seed bool) (Store, error) {
	if seed {
		err := db.AutoMigrate(&Schedule{}, &DeviceSetting{}, &Device{}, &Settings{})
		if err != nil {
			return nil, err
		}
		nightstand := &Device{
			FriendlyName: "Bedroom Nightstand",
			ID:           "3265D1FD-4FE5-4662-8AFE-C966089BCCB9",
		}
		db.Create(nightstand)
		playroomLight := &Device{
			FriendlyName: "Playroom Light",
			ID:           "3265D1FD-4FE5-4662-8AFE-C966089BCCB9",
		}
		db.Create(playroomLight)
		deskLamp1 := &Device{
			FriendlyName: "Work Desk Lamp",
			ID:           "3265D1FD-4FE5-4662-8AFE-C966089BCCA9",
		}
		db.Create(deskLamp1)
		deskLamp2 := &Device{
			FriendlyName: "Family Desk Lamp",
			ID:           "3265D1FD-5FE5-4662-8AFE-C966089BCCB9",
		}
		db.Create(deskLamp2)

		schedules := []*Schedule{
			{
				FriendlyName: "Left the Light On",
				ID:           "E0D5118D-1554-4394-93A8-EFC6C7276D0A",
				OnTime:       8 * timeofday.Hour,
				OffTime:      5*timeofday.Hour + timeofday.PM,
				DeviceSettings: []*DeviceSetting{
					{
						Device:     nightstand,
						Brightness: 100,
						Color:      "#33b73c",
					},
				},
			}, {
				FriendlyName: "Cooking",
				ID:           "31CD5DBD-E5F9-43FE-A6D3-FB7D5E07E57F",
				OnTime:       12 * timeofday.Hour,
				OffTime:      9*timeofday.Hour + timeofday.PM,
				DeviceSettings: []*DeviceSetting{
					{
						Device:     playroomLight,
						Brightness: 100,
						Color:      "#33b73c",
					},
				},
			},
			{
				ID:           "3265D1FD-4FE5-4662-8AFE-C966089BCCB0",
				FriendlyName: "Evening Puttering",
				OnTime:       8*timeofday.Hour + timeofday.PM,
				OffTime:      10*timeofday.Hour + 30*timeofday.Minute + timeofday.PM,

				DeviceSettings: []*DeviceSetting{
					{
						Device:     deskLamp1,
						Brightness: 75,
						Color:      "#33b73c",
					},
					{
						Device:     deskLamp2,
						Brightness: 75,
						Color:      "#f60080",
					},
				},
			},
		}
		for _, s := range schedules {
			db.Create(s)
		}
	}
	return &dbStore{
		database: db,
	}, nil
}
