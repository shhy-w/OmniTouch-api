package omnitouch

import (
	"context"
	"strings"
	"time"

	appSettings "git.uozi.org/uozi/cosy-example-api/settings"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/logger"
	"gorm.io/gorm"
)

type Broker struct {
	client mqtt.Client
	db     *gorm.DB
}

var broker *Broker

func InitBroker(ctx context.Context) {
	conf := appSettings.MqttSettings
	if !conf.Enabled || conf.Broker == "" {
		logger.Info("[OmniTouch MQTT] disabled, using debug simulator")
		return
	}
	db := cosy.UseDB(ctx)
	options := mqtt.NewClientOptions().
		AddBroker(conf.Broker).
		SetClientID(conf.ClientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(3 * time.Second)
	if conf.Username != "" {
		options.SetUsername(conf.Username)
		options.SetPassword(conf.Password)
	}
	b := &Broker{db: db}
	options.SetOnConnectHandler(func(client mqtt.Client) {
		logger.Info("[OmniTouch MQTT] connected")
		b.subscribe(client)
	})
	options.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logger.Error("[OmniTouch MQTT] connection lost", err)
	})
	client := mqtt.NewClient(options)
	token := client.Connect()
	token.WaitTimeout(10 * time.Second)
	if token.Error() != nil {
		logger.Error("[OmniTouch MQTT] connect failed", token.Error())
		return
	}
	b.client = client
	broker = b
}

func PublishCommand(deviceID string, payload CommandPayload) (bool, error) {
	if broker == nil || broker.client == nil || !broker.client.IsConnected() {
		return false, nil
	}
	bytes := []byte(mustJSON(payload))
	token := broker.client.Publish("term/"+deviceID+"/cmd", 1, false, bytes)
	token.WaitTimeout(5 * time.Second)
	return true, token.Error()
}

func (b *Broker) subscribe(client mqtt.Client) {
	statusToken := client.Subscribe("term/+/status", 0, b.handleStatus)
	statusToken.WaitTimeout(5 * time.Second)
	if statusToken.Error() != nil {
		logger.Error("[OmniTouch MQTT] subscribe status failed", statusToken.Error())
	}
	eventToken := client.Subscribe("term/+/event", 1, b.handleEvent)
	eventToken.WaitTimeout(5 * time.Second)
	if eventToken.Error() != nil {
		logger.Error("[OmniTouch MQTT] subscribe event failed", eventToken.Error())
	}
}

func (b *Broker) handleStatus(_ mqtt.Client, msg mqtt.Message) {
	deviceID := deviceIDFromTopic(msg.Topic())
	if deviceID == "" {
		return
	}
	if err := NewService(b.db).ApplyStatus(deviceID, msg.Payload()); err != nil {
		logger.Error("[OmniTouch MQTT] apply status failed", err)
	}
}

func (b *Broker) handleEvent(_ mqtt.Client, msg mqtt.Message) {
	deviceID := deviceIDFromTopic(msg.Topic())
	if deviceID == "" {
		return
	}
	if err := NewService(b.db).RecordEvent(deviceID, msg.Payload()); err != nil {
		logger.Error("[OmniTouch MQTT] record event failed", err)
	}
}

func deviceIDFromTopic(topic string) string {
	parts := strings.Split(topic, "/")
	if len(parts) != 3 || parts[0] != "term" {
		return ""
	}
	return parts[1]
}
