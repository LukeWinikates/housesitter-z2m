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
	Brightness      int         `json:"brightness"`
	ColorMode       string      `json:"color_mode"`
	ColorTemp       int         `json:"color_temp"`
	LevelConfig     LevelConfig `json:"level_config"`
	LinkQuality     int         `json:"linkquality"`
	State           string      `json:"state"`
	Update          Update      `json:"update"`
	UpdateAvailable bool        `json:"update_available"`
}

func OnMessage(brightness int) LightControl {
	return LightControl{
		Brightness: brightness,
		ColorMode:  "color_temp",
		ColorTemp:  370,
		LevelConfig: LevelConfig{
			OnLevel: "previous",
		},
		LinkQuality: 36,
		State:       "ON",
		Update: Update{
			InstalledVersion: 65554,
			LatestVersion:    65554,
			State:            "idle",
		},
		UpdateAvailable: false,
	}
}

func OffMessage() LightControl {
	return LightControl{
		Brightness: 0,
		ColorMode:  "color_temp",
		ColorTemp:  370,
		LevelConfig: LevelConfig{
			OnLevel: "previous",
		},
		LinkQuality: 36,
		State:       "OFF",
		Update: Update{
			InstalledVersion: 65554,
			LatestVersion:    65554,
			State:            "idle",
		},
		UpdateAvailable: false,
	}
}
