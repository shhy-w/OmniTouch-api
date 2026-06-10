package omnitouch

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) ListDevices(openid string) ([]DeviceDTO, error) {
	var devices []model.Device
	if err := s.db.Where("owner_open_id = ?", openid).Order("updated_at desc").Find(&devices).Error; err != nil {
		return nil, err
	}
	items := make([]DeviceDTO, 0, len(devices))
	for _, device := range devices {
		items = append(items, deviceDTO(device))
	}
	return items, nil
}

func (s *Service) BindDevice(openid, deviceID, name string) (DeviceDTO, error) {
	deviceID = strings.ToLower(strings.TrimSpace(deviceID))
	if deviceID == "" {
		return DeviceDTO{}, errors.New("device_id is required")
	}
	var device model.Device
	err := s.db.Where("device_id = ?", deviceID).First(&device).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		device = model.Device{
			DeviceID:    deviceID,
			OwnerOpenID: openid,
			Name:        defaultName(deviceID, name),
			Mac:         macFromDeviceID(deviceID),
			IP:          "未分配",
			Firmware:    "v1.0.0",
			PCName:      "未连接",
			StatusJSON:  mustJSON(defaultStatus(deviceID)),
			LastSeenAt:  time.Now().Unix(),
		}
		return deviceDTO(device), s.db.Create(&device).Error
	}
	if err != nil {
		return DeviceDTO{}, err
	}
	device.OwnerOpenID = openid
	if strings.TrimSpace(name) != "" {
		device.Name = strings.TrimSpace(name)
	}
	if device.StatusJSON == "" {
		device.StatusJSON = mustJSON(defaultStatus(deviceID))
	}
	if device.Mac == "" {
		device.Mac = macFromDeviceID(deviceID)
	}
	if err = s.db.Save(&device).Error; err != nil {
		return DeviceDTO{}, err
	}
	return deviceDTO(device), nil
}

func (s *Service) GetDevice(openid, deviceID string) (DeviceDTO, error) {
	device, err := s.findDevice(openid, deviceID)
	if err != nil {
		return DeviceDTO{}, err
	}
	return deviceDTO(*device), nil
}

func (s *Service) RenameDevice(openid, deviceID, name string) (DeviceDTO, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return DeviceDTO{}, errors.New("name is required")
	}
	device, err := s.findDevice(openid, deviceID)
	if err != nil {
		return DeviceDTO{}, err
	}
	device.Name = name
	if err = s.db.Save(device).Error; err != nil {
		return DeviceDTO{}, err
	}
	return deviceDTO(*device), nil
}

func (s *Service) DeleteDevice(openid, deviceID string) error {
	return s.db.Where("owner_open_id = ? AND device_id = ?", openid, deviceID).Delete(&model.Device{}).Error
}

func (s *Service) SendCommand(openid, deviceID, cmd string, params map[string]any) (*model.DeviceCommand, DeviceEventDTO, error) {
	device, err := s.findDevice(openid, deviceID)
	if err != nil {
		return nil, DeviceEventDTO{}, err
	}
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return nil, DeviceEventDTO{}, errors.New("cmd is required")
	}
	msgID := "cmd-" + uuid.NewString()
	paramsJSON := mustJSON(params)
	command := &model.DeviceCommand{
		MsgID:      msgID,
		DeviceID:   device.DeviceID,
		Cmd:        cmd,
		ParamsJSON: paramsJSON,
		State:      "sent",
		Detail:     "debug command accepted",
	}
	if err = s.db.Create(command).Error; err != nil {
		return nil, DeviceEventDTO{}, err
	}
	published, err := PublishCommand(device.DeviceID, CommandPayload{Cmd: cmd, MsgID: msgID, Params: params})
	if err != nil {
		command.State = "failed"
		command.Detail = err.Error()
		_ = s.db.Save(command).Error
		return command, DeviceEventDTO{}, err
	}
	if published {
		command.Detail = "command published to mqtt broker"
		_ = s.db.Save(command).Error
		return command, DeviceEventDTO{}, nil
	}
	event, err := s.applySimulatedCommand(device, command, params)
	return command, event, err
}

func (s *Service) SendVoice(openid, deviceID, text string) (*model.DeviceCommand, DeviceEventDTO, error) {
	cmd := voiceCommand(text)
	if cmd == "" {
		return nil, DeviceEventDTO{}, fmt.Errorf("unknown voice command: %s", text)
	}
	return s.SendCommand(openid, deviceID, cmd, nil)
}

func (s *Service) GetStatus(openid, deviceID string) (DeviceStatus, error) {
	device, err := s.findDevice(openid, deviceID)
	if err != nil {
		return DeviceStatus{}, err
	}
	return statusFromDevice(*device), nil
}

func (s *Service) ListEvents(openid, deviceID string, limit int) ([]DeviceEventDTO, error) {
	if _, err := s.findDevice(openid, deviceID); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var events []model.DeviceEvent
	if err := s.db.Where("device_id = ?", deviceID).Order("ts desc, id desc").Limit(limit).Find(&events).Error; err != nil {
		return nil, err
	}
	items := make([]DeviceEventDTO, 0, len(events))
	for _, event := range events {
		items = append(items, eventDTO(event))
	}
	return items, nil
}

func (s *Service) ApplyStatus(deviceID string, payload []byte) error {
	var status DeviceStatus
	if err := json.Unmarshal(payload, &status); err != nil {
		return err
	}
	now := time.Now().Unix()
	return s.db.Model(&model.Device{}).
		Where("device_id = ?", deviceID).
		Updates(map[string]any{
			"status_json":  string(payload),
			"last_seen_at": now,
		}).Error
}

func (s *Service) RecordEvent(deviceID string, payload []byte) error {
	var data map[string]any
	_ = json.Unmarshal(payload, &data)
	event := model.DeviceEvent{
		DeviceID:    deviceID,
		Event:       stringValue(data["event"], "event"),
		MsgID:       stringValue(data["msg_id"], ""),
		PayloadJSON: string(payload),
		Ts:          int64Value(data["ts"], time.Now().Unix()),
	}
	if event.Event == "" {
		event.Event = "event"
	}
	return s.db.Create(&event).Error
}

func (s *Service) findDevice(openid, deviceID string) (*model.Device, error) {
	var device model.Device
	err := s.db.Where("owner_open_id = ? AND device_id = ?", openid, strings.ToLower(deviceID)).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (s *Service) applySimulatedCommand(device *model.Device, command *model.DeviceCommand, params map[string]any) (DeviceEventDTO, error) {
	status := statusFromDevice(*device)
	switch command.Cmd {
	case "lock":
		status.Locked = true
	case "unlock":
		status.Locked = false
	case "set_mode":
		if mode, ok := params["mode"].(string); ok && mode != "" {
			status.Mode = mode
		}
	case "laser_on":
		status.LaserOn = true
	case "laser_off":
		status.LaserOn = false
	case "set_backlight":
		status.Backlight = clampInt(anyInt(params["value"], status.Backlight), 0, 100)
	case "set_backlight_auto":
		if value, ok := params["enable"].(bool); ok {
			status.BacklightAuto = value
		}
	case "zoom_in":
		status.Zoom = math.Min(status.Zoom+0.1*float64(anyInt(params["level"], 1)), 2)
	case "zoom_out":
		status.Zoom = math.Max(status.Zoom-0.1*float64(anyInt(params["level"], 1)), 0.5)
	case "next_page":
		if status.CurrentPage < status.TotalPages {
			status.CurrentPage++
		}
	case "prev_page":
		if status.CurrentPage > 1 {
			status.CurrentPage--
		}
	}
	status.Online = true
	status.LastSeen = "刚刚"
	device.StatusJSON = mustJSON(status)
	device.LastSeenAt = time.Now().Unix()
	if err := s.db.Save(device).Error; err != nil {
		return DeviceEventDTO{}, err
	}
	event := DeviceEventDTO{
		Event:  "cmd_done",
		MsgID:  command.MsgID,
		Detail: command.Cmd + " executed",
		Ts:     time.Now().Unix(),
	}
	record := model.DeviceEvent{
		DeviceID:    device.DeviceID,
		Event:       event.Event,
		MsgID:       event.MsgID,
		PayloadJSON: mustJSON(event),
		Ts:          event.Ts,
	}
	if err := s.db.Create(&record).Error; err != nil {
		return DeviceEventDTO{}, err
	}
	command.State = "done"
	command.Detail = event.Detail
	_ = s.db.Save(command).Error
	return event, nil
}

func deviceDTO(device model.Device) DeviceDTO {
	return DeviceDTO{ID: device.DeviceID, Name: device.Name, Status: statusFromDevice(device)}
}

func statusFromDevice(device model.Device) DeviceStatus {
	status := defaultStatus(device.DeviceID)
	if device.StatusJSON != "" {
		_ = json.Unmarshal([]byte(device.StatusJSON), &status)
	}
	status.Mac = firstNonEmpty(status.Mac, device.Mac)
	status.IP = firstNonEmpty(status.IP, device.IP)
	status.Firmware = firstNonEmpty(status.Firmware, device.Firmware)
	status.PCName = firstNonEmpty(status.PCName, device.PCName)
	if device.LastSeenAt > 0 && time.Now().Unix()-device.LastSeenAt > 30 {
		status.LastSeen = formatLastSeen(device.LastSeenAt)
	}
	return status
}

func eventDTO(event model.DeviceEvent) DeviceEventDTO {
	dto := DeviceEventDTO{Event: event.Event, MsgID: event.MsgID, Ts: event.Ts}
	_ = json.Unmarshal([]byte(event.PayloadJSON), &dto)
	if dto.Event == "" {
		dto.Event = event.Event
	}
	if dto.Ts == 0 {
		dto.Ts = event.Ts
	}
	return dto
}

func defaultStatus(deviceID string) DeviceStatus {
	return DeviceStatus{
		Online:        true,
		Locked:        false,
		Mode:          "mqtt",
		Battery:       85,
		Charging:      false,
		RSSI:          -45,
		CurrentPage:   1,
		TotalPages:    20,
		PCName:        "未连接",
		Firmware:      "v1.0.0",
		IP:            "未分配",
		Mac:           macFromDeviceID(deviceID),
		LaserOn:       false,
		Backlight:     75,
		BacklightAuto: true,
		Zoom:          1,
		LastSeen:      "刚刚",
	}
}

func defaultName(deviceID, name string) string {
	if strings.TrimSpace(name) != "" {
		return strings.TrimSpace(name)
	}
	suffix := strings.ToUpper(deviceID)
	if len(suffix) > 6 {
		suffix = suffix[len(suffix)-6:]
	}
	return "SmartTerm_" + suffix
}

func macFromDeviceID(deviceID string) string {
	suffix := strings.ToUpper(strings.TrimPrefix(deviceID, "esp32-"))
	if len(suffix) < 6 {
		return ""
	}
	suffix = suffix[len(suffix)-6:]
	return fmt.Sprintf("24:6F:28:%s:%s:%s", suffix[0:2], suffix[2:4], suffix[4:6])
}

func voiceCommand(text string) string {
	switch strings.TrimSpace(text) {
	case "下一页", "下翻":
		return "next_page"
	case "上一页", "上翻":
		return "prev_page"
	case "放大":
		return "zoom_in"
	case "缩小":
		return "zoom_out"
	case "锁定":
		return "lock"
	case "解锁":
		return "unlock"
	default:
		return ""
	}
}

func mustJSON(value any) string {
	bytes, _ := json.Marshal(value)
	return string(bytes)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func anyInt(value any, fallback int) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		i, _ := v.Int64()
		return int(i)
	default:
		return fallback
	}
}

func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func stringValue(value any, fallback string) string {
	if text, ok := value.(string); ok {
		return text
	}
	return fallback
}

func int64Value(value any, fallback int64) int64 {
	switch v := value.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	default:
		return fallback
	}
}

func formatLastSeen(ts int64) string {
	diff := time.Now().Unix() - ts
	if diff < 60 {
		return "刚刚"
	}
	if diff < 3600 {
		return fmt.Sprintf("%d 分钟前", diff/60)
	}
	return fmt.Sprintf("%d 小时前", diff/3600)
}
