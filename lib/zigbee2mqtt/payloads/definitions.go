package payloads

import "encoding/json"

type Definition struct {
	Description string `json:"description"`
	Exposes     []struct {
		Access      int      `json:"access,omitempty"`
		Category    string   `json:"category,omitempty"`
		Description string   `json:"description,omitempty"`
		Label       string   `json:"label,omitempty"`
		Name        string   `json:"name,omitempty"`
		Property    string   `json:"property,omitempty"`
		Type        string   `json:"type"`
		Values      []string `json:"values,omitempty"`
		Unit        string   `json:"unit,omitempty"`
		ValueMax    int      `json:"value_max,omitempty"`
		ValueMin    int      `json:"value_min,omitempty"`
		Features    []struct {
			Access      int         `json:"access"`
			Description string      `json:"description"`
			Label       string      `json:"label"`
			Name        string      `json:"name"`
			Property    string      `json:"property"`
			Type        string      `json:"type"`
			ValueOff    interface{} `json:"value_off,omitempty"`
			ValueOn     interface{} `json:"value_on,omitempty"`
			ValueToggle string      `json:"value_toggle,omitempty"`
			ValueMax    int         `json:"value_max,omitempty"`
			ValueMin    int         `json:"value_min,omitempty"`
			Presets     []struct {
				Description string `json:"description"`
				Name        string `json:"name"`
				Value       int    `json:"value"`
			} `json:"presets,omitempty"`
			Unit     string `json:"unit,omitempty"`
			Features []struct {
				Access   int    `json:"access"`
				Label    string `json:"label"`
				Name     string `json:"name"`
				Property string `json:"property"`
				Type     string `json:"type"`
			} `json:"features,omitempty"`
		} `json:"features,omitempty"`
	} `json:"exposes"`
	Model   string `json:"model"`
	Options []struct {
		Access      int    `json:"access"`
		Description string `json:"description"`
		Label       string `json:"label"`
		Name        string `json:"name"`
		Property    string `json:"property"`
		Type        string `json:"type"`
		ValueMax    int    `json:"value_max,omitempty"`
		ValueMin    int    `json:"value_min,omitempty"`
		Features    []struct {
			Access      int    `json:"access"`
			Description string `json:"description"`
			Label       string `json:"label"`
			Name        string `json:"name"`
			Property    string `json:"property"`
			Type        string `json:"type"`
			ValueMin    int    `json:"value_min"`
			Unit        string `json:"unit,omitempty"`
		} `json:"features,omitempty"`
		ValueOff bool `json:"value_off,omitempty"`
		ValueOn  bool `json:"value_on,omitempty"`
	} `json:"options"`
	SupportsOta bool   `json:"supports_ota"`
	Vendor      string `json:"vendor"`
}

type Endpoint struct {
	Bindings []struct {
		Cluster string `json:"cluster"`
		Target  struct {
			Endpoint    int    `json:"endpoint"`
			IeeeAddress string `json:"ieee_address"`
			Type        string `json:"type"`
		} `json:"target"`
	} `json:"bindings"`
	Clusters struct {
		Input  []string `json:"input"`
		Output []string `json:"output"`
	} `json:"clusters"`
	ConfiguredReportings []struct {
		Attribute             string `json:"attribute"`
		Cluster               string `json:"cluster"`
		MaximumReportInterval int    `json:"maximum_report_interval"`
		MinimumReportInterval int    `json:"minimum_report_interval"`
		ReportableChange      int    `json:"reportable_change"`
	} `json:"configured_reportings"`
	Scenes []interface{} `json:"scenes"`
}

type MessagePayload struct {
	Definition         *Definition         `json:"definition"`
	Disabled           bool                `json:"disabled"`
	Endpoints          map[string]Endpoint `json:"endpoints"`
	FriendlyName       string              `json:"friendly_name"`
	IEEEAddress        string              `json:"ieee_address"`
	InterviewCompleted bool                `json:"interview_completed"`
	Interviewing       bool                `json:"interviewing"`
	NetworkAddress     int                 `json:"network_address"`
	Supported          bool                `json:"supported"`
	Type               string              `json:"type"`
	DateCode           string              `json:"date_code,omitempty"`
	Manufacturer       string              `json:"manufacturer,omitempty"`
	ModelId            string              `json:"model_id,omitempty"`
	PowerSource        string              `json:"power_source,omitempty"`
	SoftwareBuildId    string              `json:"software_build_id,omitempty"`
}

func Parse(payload []byte) ([]MessagePayload, error) {
	var parseResult []MessagePayload
	err := json.Unmarshal(payload, &parseResult)
	return parseResult, err
}
