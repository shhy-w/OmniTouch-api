package model

type MiniappUser struct {
	Model
	OpenID      string `json:"openid" gorm:"type:varchar(128);uniqueIndex"`
	UnionID     string `json:"unionid,omitempty" gorm:"type:varchar(128);index"`
	SessionKey  string `json:"-" gorm:"type:varchar(255)"`
	LastLoginAt int64  `json:"last_login_at,omitempty"`
}

type Device struct {
	Model
	DeviceID    string `json:"id" gorm:"type:varchar(64);uniqueIndex"`
	OwnerOpenID string `json:"owner_openid,omitempty" gorm:"type:varchar(128);index"`
	Name        string `json:"name" gorm:"type:varchar(128)"`
	Mac         string `json:"mac,omitempty" gorm:"type:varchar(64)"`
	IP          string `json:"ip,omitempty" gorm:"type:varchar(64)"`
	Firmware    string `json:"firmware,omitempty" gorm:"type:varchar(64)"`
	PCName      string `json:"pc_name,omitempty" gorm:"type:varchar(128)"`
	StatusJSON  string `json:"-" gorm:"type:json"`
	LastSeenAt  int64  `json:"last_seen_at,omitempty" gorm:"index"`
}

type DeviceCommand struct {
	Model
	MsgID      string `json:"msg_id" gorm:"type:varchar(128);uniqueIndex"`
	DeviceID   string `json:"device_id" gorm:"type:varchar(64);index"`
	Cmd        string `json:"cmd" gorm:"type:varchar(64);index"`
	ParamsJSON string `json:"-" gorm:"type:json"`
	State      string `json:"state" gorm:"type:varchar(32);index"`
	Detail     string `json:"detail,omitempty" gorm:"type:varchar(255)"`
}

type DeviceEvent struct {
	Model
	DeviceID    string `json:"device_id" gorm:"type:varchar(64);index"`
	Event       string `json:"event" gorm:"type:varchar(64);index"`
	MsgID       string `json:"msg_id,omitempty" gorm:"type:varchar(128);index"`
	PayloadJSON string `json:"-" gorm:"type:json"`
	Ts          int64  `json:"ts" gorm:"index"`
}
