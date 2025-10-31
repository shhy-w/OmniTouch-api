package model

import (
	"github.com/uozi-tech/cosy/sonyflake"
	"gorm.io/gorm"
)

type BurnedHistory struct {
	Model
	EnKey          string        `json:"en_key" cosy:"add:required;update:omitempty;list:fussy"`
	FirmwareFileId uint64        `json:"firmware_file_id" cosy:"add:required;update:omitempty;list:fussy"`
	// FirmwareFile   *FirmwareFile `json:"firmware_file" cosy:"update:omitempty;list:preload;item:preload;"`//该 model 已被删除
	ClientIp       string        `json:"client_ip" cosy:"add:required;update:omitempty;list:fussy"`
	Status         BurnStatus    `json:"status" cosy:"add:required;update:omitempty;list:fussy"`
	Remark         string        `json:"remark" cosy:"list:fussy"`
}

// BeforeCreate 是一个钩子函数，在创建 BurnedHistory 对象之前被调用
func (bh *BurnedHistory) BeforeCreate(_ *gorm.DB) error {
	if bh.ID == 0 {
		bh.ID = sonyflake.NextID()
	}
	return nil
}

type BurnStatus string

const (
	BurnStatusSuccess  BurnStatus = "success"
	BurnStatusFailed   BurnStatus = "failed"
	Burning            BurnStatus = "burning"
	BurnStatusCanceled BurnStatus = "canceled"
	Burned             BurnStatus = "burned"
)
