package queries

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/timeofday"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceSettingCreate struct {
	ID         string
	Brightness int
	Color      string
}

type ScheduleCreateTemplate struct {
	Name    string
	OnTime  timeofday.SecondsInDay
	OffTime timeofday.SecondsInDay
	Devices []*DeviceSettingCreate
}

func CreateSchedule(db *gorm.DB, request ScheduleCreateTemplate) error {
	newSchedule := &database.Schedule{
		ID:           uuid.NewString(),
		FriendlyName: request.Name,
		OffTime:      request.OffTime,
		OnTime:       request.OnTime,
	}

	for _, requestedDevice := range request.Devices {
		device := &database.Device{}
		err := db.Find(device, "id = ?", requestedDevice.ID).Error
		if err != nil {
			return fmt.Errorf("device not found for ID %s", requestedDevice.ID)
		}
		err = db.Model(newSchedule).Association("DeviceSettings").Append(&database.DeviceSetting{
			Device:     device,
			Brightness: uint8(requestedDevice.Brightness),
			Color:      requestedDevice.Color,
			DeviceID:   device.ID,
		})
		if err != nil {
			return err
		}
	}
	return db.Create(newSchedule).Error
}
