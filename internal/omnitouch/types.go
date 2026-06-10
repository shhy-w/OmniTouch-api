package omnitouch

type DeviceStatus struct {
	Online        bool    `json:"online"`
	Locked        bool    `json:"locked"`
	Mode          string  `json:"mode"`
	Battery       int     `json:"battery"`
	Charging      bool    `json:"charging"`
	RSSI          int     `json:"rssi"`
	CurrentPage   int     `json:"currentPage"`
	TotalPages    int     `json:"totalPages"`
	PCName        string  `json:"pcName"`
	Firmware      string  `json:"firmware"`
	IP            string  `json:"ip"`
	Mac           string  `json:"mac"`
	LaserOn       bool    `json:"laserOn"`
	Backlight     int     `json:"backlight"`
	BacklightAuto bool    `json:"backlightAuto"`
	Zoom          float64 `json:"zoom"`
	LastSeen      string  `json:"lastSeen"`
}

type DeviceDTO struct {
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Status DeviceStatus `json:"status"`
}

type DeviceEventDTO struct {
	Event      string      `json:"event"`
	MsgID      string      `json:"msg_id,omitempty"`
	Detail     string      `json:"detail,omitempty"`
	Reason     string      `json:"reason,omitempty"`
	Command    string      `json:"command,omitempty"`
	GestureID  string      `json:"gesture_id,omitempty"`
	Confidence float64     `json:"confidence,omitempty"`
	Ts         int64       `json:"ts"`
	Payload    interface{} `json:"payload,omitempty"`
}

type CommandPayload struct {
	Cmd    string         `json:"cmd"`
	MsgID  string         `json:"msg_id"`
	Params map[string]any `json:"params,omitempty"`
}
