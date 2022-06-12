package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
)

var key = []byte("passphrasewhichneedstobe32bytes!")
var IV = []byte("1234567812345678")

//the use of generated key file is preferable but hey - this is life
//if you got this far and read it than you are one of the rare courageous and bold
//FIX it if you need to - but only if you really need to
const (
	keyFile       = "aes.key"
	encryptedFile = "aes.enc"
)

func readKey(filename string) ([]byte, error) {
	key, err := ioutil.ReadFile(filename)
	if err != nil {
		return key, err
	}
	block, _ := pem.Decode(key)
	return block.Bytes, nil
}

func createKey() []byte {
	genkey := make([]byte, 16)
	_, err := rand.Read(genkey)
	if err != nil {
		log.Fatalf("Failed to read new random key: %s", err)
	}
	return genkey
}

func saveKey(filename string, key []byte) {
	block := &pem.Block{
		Type:  "AES KEY",
		Bytes: key,
	}
	err := ioutil.WriteFile(filename, pem.EncodeToMemory(block), 0644)
	if err != nil {
		log.Fatalf("Failed in saving key to %s: %s", filename, err)
	}
}

func aesKey() []byte {
	file := fmt.Sprintf(keyFile)
	key, err := readKey(file)
	if err != nil {
		log.Println("Creating a new AES key")
		key = createKey()
		saveKey(file, key)
	}
	return key
}

func createCipherFromFile() cipher.Block {
	c, err := aes.NewCipher(aesKey())
	if err != nil {
		log.Fatalf("Failed to create the AES cipher: %s", err)
	}
	return c
}

func createCipherFromString() cipher.Block {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Failed to create the AES cipher: %s", err)
	}
	return c
}

func Encryption(plainText string) string{
	bytes := []byte(plainText)
	blockCipher := createCipherFromString()
	stream := cipher.NewCTR(blockCipher, IV)
	stream.XORKeyStream(bytes, bytes)
	return string(bytes)
}

func Decryption(encryptedMessage string) string {
	blockCipher := createCipherFromString()
	stream := cipher.NewCTR(blockCipher, IV)
	bytes := []byte(encryptedMessage)
	stream.XORKeyStream(bytes, bytes)
	return string(bytes)
}

