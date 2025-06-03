package queries

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"fmt"
	"gorm.io/gorm"
)

func UpdateSchedule(db *gorm.DB, existingSchedule *database.Schedule, request ScheduleCreateTemplate) error {
	existingSchedule.FriendlyName = request.Name
	existingSchedule.OffTime = request.OffTime
	existingSchedule.OnTime = request.OnTime

	if err := db.Save(existingSchedule).Error; err != nil {
		return err
	}

	if err := db.Model(existingSchedule).Association("DeviceSettings").Clear(); err != nil {
		return err
	}

	for _, requestedDevice := range request.Devices {
		device := &database.Device{}
		err := db.Find(device, "id = ?", requestedDevice.ID).Error
		if err != nil {
			return fmt.Errorf("device not found for ID %s", requestedDevice.ID)
		}
		err = db.Model(existingSchedule).Association("DeviceSettings").Append(&database.DeviceSetting{
			Device:     device,
			Brightness: uint8(requestedDevice.Brightness),
			Color:      requestedDevice.Color,
			DeviceID:   device.ID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
