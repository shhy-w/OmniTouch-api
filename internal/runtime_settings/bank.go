package runtime_settings

import (
	"encoding/json"
	"github.com/uozi-tech/cosy/logger"
)

type BankInstruction struct {
	BankName      string `json:"bank_name"`
	SwiftCode     string `json:"swift_code"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
}

func GetBankInstructionSettings() (a *BankInstruction) {
	a = &BankInstruction{}
	err := Get("bank", a)
	if err != nil {
		logger.Error(err)
	}
	return
}

func UpdateBankInstructionSettings(target *BankInstruction) {
	bytes, err := json.Marshal(target)
	if err != nil {
		logger.Error(err)
		return
	}
	Set("bank", string(bytes))
}
