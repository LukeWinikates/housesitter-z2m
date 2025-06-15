package devices

type Update struct {
	InstalledVersion int    `json:"installed_version"`
	LatestVersion    int    `json:"latest_version"`
	State            string `json:"state"`
}

type LevelConfig struct {
	OnLevel string `json:"on_level"`
}

type LightControl struct {
	Brightness      *int         `json:"brightness,omitempty"`
	ColorMode       string      `json:"color_mode,omitempty"`
	ColorTemp       int         `json:"color_temp,omitempty"`
	LevelConfig     LevelConfig `json:"level_config,omitempty"`
	LinkQuality     int         `json:"linkquality,omitempty"`
	State           string      `json:"state,omitempty"`
	Update          Update      `json:"update,omitempty"`
	UpdateAvailable bool        `json:"update_available,omitempty"`
}

func OnMessage(brightness int) LightControl {
	return LightControl{
		Brightness: &brightness,
		State:       "ON",
	}
}

func OffMessage() LightControl {
	return LightControl{
		State:       "OFF",
	}
}
