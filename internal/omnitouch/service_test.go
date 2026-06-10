package omnitouch

import (
	"encoding/json"
	"testing"

	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testService(t *testing.T) *Service {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.Device{},
		&model.DeviceCommand{},
		&model.DeviceEvent{},
	))
	return NewService(db)
}

func TestBindDeviceAndCommandSimulation(t *testing.T) {
	service := testService(t)
	device, err := service.BindDevice("openid-1", "esp32-cb8824", "")
	require.NoError(t, err)
	require.Equal(t, "esp32-cb8824", device.ID)
	require.Equal(t, "SmartTerm_CB8824", device.Name)

	command, event, err := service.SendCommand("openid-1", "esp32-cb8824", "next_page", nil)
	require.NoError(t, err)
	require.Equal(t, "done", command.State)
	require.Equal(t, "cmd_done", event.Event)

	status, err := service.GetStatus("openid-1", "esp32-cb8824")
	require.NoError(t, err)
	require.Equal(t, 2, status.CurrentPage)
}

func TestApplyStatusAndRecordEvent(t *testing.T) {
	service := testService(t)
	_, err := service.BindDevice("openid-1", "esp32-cb8824", "")
	require.NoError(t, err)

	payload, _ := json.Marshal(DeviceStatus{
		Online:      true,
		Mode:        "gesture",
		Battery:     66,
		RSSI:        -52,
		CurrentPage: 4,
		TotalPages:  10,
	})
	require.NoError(t, service.ApplyStatus("esp32-cb8824", payload))

	status, err := service.GetStatus("openid-1", "esp32-cb8824")
	require.NoError(t, err)
	require.Equal(t, "gesture", status.Mode)
	require.Equal(t, 66, status.Battery)

	eventPayload := []byte(`{"event":"wake","ts":1717849320}`)
	require.NoError(t, service.RecordEvent("esp32-cb8824", eventPayload))
	events, err := service.ListEvents("openid-1", "esp32-cb8824", 20)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "wake", events[0].Event)
}
