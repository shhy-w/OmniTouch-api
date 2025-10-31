package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	wyKey = "61696f74636f6d6d"
)

// hexToBinary 将单个十六进制字符转换为4位二进制字符串
func hexToBinary(hexChar rune) string {
	switch hexChar {
	case '0':
		return "0000"
	case '1':
		return "0001"
	case '2':
		return "0010"
	case '3':
		return "0011"
	case '4':
		return "0100"
	case '5':
		return "0101"
	case '6':
		return "0110"
	case '7':
		return "0111"
	case '8':
		return "1000"
	case '9':
		return "1001"
	case 'a', 'A':
		return "1010"
	case 'b', 'B':
		return "1011"
	case 'c', 'C':
		return "1100"
	case 'd', 'D':
		return "1101"
	case 'e', 'E':
		return "1110"
	case 'f', 'F':
		return "1111"
	default:
		// TODO: 这里需要处理确认非法字符要怎么办
		return "0000"
	}
}

// generateAESKey 根据 enKey 生成 AES 密钥
func GenerateAESKey(enKey string) ([]byte, error) {
	// 确保 enKey 是 64 字符长度
	if len(enKey) != 64 {
		return nil, fmt.Errorf("enKey must be 64 characters long, got %d", len(enKey))
	}

	// 将 enkey 的奇偶位反转（相邻字符交换位置）
	var swappedKey strings.Builder
	runes := []rune(enKey)
	for i := 0; i < len(runes); i += 2 {
		if i+1 < len(runes) {
			// 交换相邻的两个字符
			swappedKey.WriteRune(runes[i+1])
			swappedKey.WriteRune(runes[i])
		} else {
			// 如果是奇数长度（虽然这里保证是64位），最后一个字符不变
			swappedKey.WriteRune(runes[i])
		}
	}
	swappedEnKey := swappedKey.String()

	// 将最后 16 位替换成 wyKey
	modifiedKey := swappedEnKey[:48] + wyKey

	// 将 64 个十六进制字符转换为 256 位二进制字符串
	var binaryStr strings.Builder
	for _, char := range modifiedKey {
		binaryStr.WriteString(hexToBinary(char))
	}

	// 将 256 位二进制字符串转换为 32 字节的密钥
	binaryString := binaryStr.String()
	if len(binaryString) != 256 {
		return nil, fmt.Errorf("binary string length should be 256, got %d", len(binaryString))
	}

	key := make([]byte, 32) // AES-256 需要 32 字节密钥
	for i := 0; i < 32; i++ {
		byteStr := binaryString[i*8 : (i+1)*8]
		byteVal, err := strconv.ParseUint(byteStr, 2, 8)
		if err != nil {
			return nil, fmt.Errorf("failed to parse binary byte: %v", err)
		}
		key[i] = byte(byteVal)
	}

	return key, nil
}

// encryptFile 使用 AES-GCM 模式加密文件
func EncryptFile(inputPath, outputPath string, key []byte) error {
	// 读取原始文件
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %v", err)
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %v", err)
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %v", err)
	}

	// 加密数据
	cipherText := gcm.Seal(nonce, nonce, plaintext, nil)

	// 写入加密后的文件
	err = os.WriteFile(outputPath, cipherText, 0644)
	if err != nil {
		return fmt.Errorf("failed to write encrypted file: %v", err)
	}

	return nil
}

// DecryptFile 使用 AES-GCM 模式解密文件
func DecryptFile(inputPath, outputPath string, key []byte) error {
	// 读取加密文件
	cipherText, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %v", err)
	}

	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %v", err)
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %v", err)
	}

	// 检查加密数据长度
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return fmt.Errorf("encrypted data too short")
	}

	// 提取 nonce 和加密数据
	nonce := cipherText[:nonceSize]
	encryptedData := cipherText[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %v", err)
	}

	// 写入解密后的文件
	err = os.WriteFile(outputPath, plaintext, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decrypted file: %v", err)
	}

	return nil
}

// DecryptFirmwareWithKey 根据 enKey 解密固件文件的便捷函数
func DecryptFirmwareWithKey(enKey, inputPath, outputPath string) error {
	// 生成 AES 密钥
	aesKey, err := GenerateAESKey(enKey)
	if err != nil {
		return fmt.Errorf("failed to generate AES key: %v", err)
	}

	// 解密文件
	err = DecryptFile(inputPath, outputPath, aesKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt file: %v", err)
	}

	return nil
}
