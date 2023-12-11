package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/Lockps/Forres-release-version/cmd/database"
)

func Encoder(permission int, key []byte) {
	filePath := database.GetLocation(permission) + ".db"

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		fmt.Println("Error generating key:", err)
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating AES cipher:", err)
		return
	}

	ciphertext := make([]byte, aes.BlockSize+len(fileContent))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("Error generating IV:", err)
		return
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], fileContent)

	err = os.WriteFile(filePath, ciphertext, 0644)
	if err != nil {
		fmt.Println("Error writing encrypted content to file:", err)
		return
	}

	fmt.Println("File content encrypted using AES.")
}

func DeCoder(permission int, key []byte) {
	filePath := database.GetLocation(permission) + ".db"

	encryptedContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading encrypted file:", err)
		return
	}

	iv := encryptedContent[:aes.BlockSize]
	encryptedContent = encryptedContent[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating AES cipher:", err)
		return
	}

	plaintext := make([]byte, len(encryptedContent))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, encryptedContent)

	decryptedFilePath := "path_to_save_decrypted_file.db"
	err = os.WriteFile(decryptedFilePath, plaintext, 0644)
	if err != nil {
		fmt.Println("Error writing decrypted content to file:", err)
		return
	}

	fmt.Println("File content decrypted using AES.")
}

func EncoderBase64(permission int) {
	filePath := database.GetLocation(permission) + ".db"

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	encodedContent := base64.StdEncoding.EncodeToString(fileContent)

	err = os.WriteFile(filePath, []byte(encodedContent), 0644)
	if err != nil {
		fmt.Println("Error writing encoded content to file:", err)
		return
	}

	fmt.Println("File content encoded to Base64.")
}

func DecoderBase64(permission int) {
	filePath := database.GetLocation(permission) + ".db"

	encodedContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	decodedContent, err := base64.StdEncoding.DecodeString(string(encodedContent))
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	err = os.WriteFile(filePath, decodedContent, 0644)
	if err != nil {
		fmt.Println("Error writing decoded content to file:", err)
		return
	}

	fmt.Println("File content decoded from Base64.")
}
