package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
)

const (
	saltSize  = 16
	iteration = 60000_0
	keySize   = 32
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n1. Зашифровать\n2. Расшифровать\nВыбор: ")
		if !sc.Scan() {
			return
		}
		switch sc.Text() {
		case "1":
			text := ask(sc, "Текст: ")
			password := ask(sc, "Ключ: ")
			res, err := encrypt(text, password)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			fmt.Println("Зашифровано:", res)
		case "2":
			encoded := ask(sc, "Шифротекст (base64): ")
			password := ask(sc, "Ключ: ")
			res, err := decrypt(encoded, password)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			fmt.Println("Расшифровано:", res)
		default:
			fmt.Println("Нет такого пункта")
		}
	}
}

func ask(sc *bufio.Scanner, label string) string {
	fmt.Print(label)
	sc.Scan()
	return sc.Text()
}

func deriveKey(password string, salt []byte) ([]byte, error) {
	return pbkdf2.Key(sha256.New, password, salt, iteration, keySize)
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
func encrypt(plaintext, password string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key, err := deriveKey(password, salt)
	if err != nil {
		return "", err
	}
	gcm, err := newGCM(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	out := append(salt, nonce...)
	out = gcm.Seal(out, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(out), nil
}

func decrypt(encoded, password string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", errors.New("это не base64")
	}
	if len(data) < saltSize+12 {
		return "", errors.New("слишком короткий")
	}
	salt, rest := data[:saltSize], data[saltSize:]
	key, err := deriveKey(password, salt)
	if err != nil {
		return "", err
	}
	gcm, err := newGCM(key)
	if err != nil {
		return "", err
	}
	nonce, ct := rest[:gcm.NonceSize()], rest[gcm.NonceSize():]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", errors.New("неверный ключ или данные повреждены")
	}
	return string(pt), nil
}
