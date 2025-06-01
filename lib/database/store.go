package database

type Store interface {
	Find(id string) (*Schedule, error)
	All() []*Schedule
	SaveChanges(id string, s *Schedule) error
	FindScheduleDeviceSettings(scheduleId, deviceId string) (*DeviceSetting, error)
	SaveDeviceSettingChanges(scheduleId string, deviceId string, settings *DeviceSetting) error
	Add(schedule *Schedule) error
	AddDeviceSettings(scheduleID string, settings *DeviceSetting) error
}
