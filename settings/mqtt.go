package settings

type Mqtt struct {
	Broker   string
	Username string
	Password string
	ClientID string
	Enabled  bool
}

var MqttSettings = &Mqtt{
	ClientID: "omnitouch-api",
	Enabled:  false,
}
