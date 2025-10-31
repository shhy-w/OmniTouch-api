package aes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptDecryptFlow(t *testing.T) {
	// 测试用的 enKey (64字符)
	enKey := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	originalFile := filepath.Join(tempDir, "test_firmware.bin")
	encryptedFile := filepath.Join(tempDir, "encrypted_firmware.bin")
	decryptedFile := filepath.Join(tempDir, "decrypted_firmware.bin")

	// 写入测试数据
	testData := []byte("This is a test firmware file content for encryption testing.")
	err := os.WriteFile(originalFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 生成 AES 密钥
	aesKey, err := GenerateAESKey(enKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	// 加密文件
	err = EncryptFile(originalFile, encryptedFile, aesKey)
	if err != nil {
		t.Fatalf("Failed to encrypt file: %v", err)
	}

	// 验证加密文件存在
	if _, err := os.Stat(encryptedFile); os.IsNotExist(err) {
		t.Fatalf("Encrypted file does not exist")
	}

	// 解密文件
	err = DecryptFile(encryptedFile, decryptedFile, aesKey)
	if err != nil {
		t.Fatalf("Failed to decrypt file: %v", err)
	}

	// 读取解密后的数据
	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("Failed to read decrypted file: %v", err)
	}

	// 验证解密后的数据与原始数据相同
	if string(decryptedData) != string(testData) {
		t.Fatalf("Decrypted data does not match original data")
	}

	t.Logf("Encryption and decryption test passed successfully")
}

func TestDecryptFirmwareWithKey(t *testing.T) {
	// 测试用的 enKey (64字符)
	enKey := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"

	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	originalFile := filepath.Join(tempDir, "firmware.bin")
	encryptedFile := filepath.Join(tempDir, "encrypted_firmware.bin")
	decryptedFile := filepath.Join(tempDir, "decrypted_firmware.bin")

	// 写入测试数据
	testData := []byte("Ffbiusafgaufauyrewyrawyeryw8rywyfsadhfdsuhfsiirmware data with encryption test")
	err := os.WriteFile(originalFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 先加密文件 (模拟服务器端的加密过程)
	aesKey, err := GenerateAESKey(enKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	err = EncryptFile(originalFile, encryptedFile, aesKey)
	if err != nil {
		t.Fatalf("Failed to encrypt file: %v", err)
	}

	// 使用便捷函数解密 (这是用户会使用的函数)
	err = DecryptFirmwareWithKey(enKey, encryptedFile, decryptedFile)
	if err != nil {
		t.Fatalf("Failed to decrypt firmware with key: %v", err)
	}

	// 验证解密结果
	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("Failed to read decrypted file: %v", err)
	}
	t.Log("decryptedData:----------------------------------")
	t.Log(string(decryptedData))
	t.Log("testData:----------------------------------")
	t.Log(string(testData))
	if string(decryptedData) != string(testData) {
		t.Fatalf("Decrypted data does not match original data")
	}

	t.Logf("DecryptFirmwareWithKey test passed successfully")
}
