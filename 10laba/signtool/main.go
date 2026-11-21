package main

import (
	"bufio"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

type SignedMessage struct {
	From      string `json:"from"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func generateKeys(name string) error {

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}
	privPEM := pem.EncodeToMemory(
		&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if err := os.WriteFile(name+".key", privPEM, 0600); err != nil {
		return err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return err
	}
	pubPEM := pem.EncodeToMemory(
		&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	if err := os.WriteFile(name+".pub", pubPEM, 0644); err != nil {
		return err
	}
	return nil
}

func loadPublicKey(path string) (ed25519.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("в файле нет PEM")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("не тот ключ")
	}
	return pub, nil
}
func loadPrivateKey(path string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("в файле нет PEM")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("не тот ключ")
	}
	return priv, nil
}

func signMessage(from, message string) error {
	priv, err := loadPrivateKey(from + ".key")
	if err != nil {
		return err
	}
	sig := ed25519.Sign(priv, []byte(message))
	msg := SignedMessage{From: from, Message: message, Signature: base64.StdEncoding.EncodeToString(sig)}
	data, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(from+".json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
func verifyLetter(letterPath, pubPath string) (bool, error) {
	data, err := os.ReadFile(letterPath)
	if err != nil {
		return false, err
	}
	var msg SignedMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return false, err
	}
	pub, err := loadPublicKey(pubPath)
	if err != nil {
		return false, err
	}

	sig, err := base64.StdEncoding.DecodeString(msg.Signature)
	if err != nil {
		return false, err
	}
	return ed25519.Verify(pub, []byte(msg.Message), sig), nil

}
func main() {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n1. Создать пару ключей\n2. Подписать сообщение\n3. Проверить письмо\nВыбор: ")
		if !sc.Scan() {
			return
		}
		switch sc.Text() {
		case "1":
			name := ask(sc, "Имя владельца: ")
			if err := generateKeys(name); err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			fmt.Printf("Созданы %s.key и %s.pub\n", name, name)
		case "2":
			from := ask(sc, "Отправитель: ")
			message := ask(sc, "Сообщение: ")
			if err := signMessage(from, message); err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			fmt.Printf("Письмо подписано и сохранено в %s.json\n", from)
		case "3":
			letter := ask(sc, "Файл письма: ")
			pub := ask(sc, "Файл открытого ключа: ")
			ok, err := verifyLetter(letter, pub)
			if err != nil {
				fmt.Println("Не удалось проверить:", err)
				continue
			}
			if ok {
				fmt.Println("Подпись верна")
			} else {
				fmt.Println("Подпись неверна")
			}
		default:
			fmt.Println("Нет такого пункта")
			return
		}
	}
}

func ask(sc *bufio.Scanner, label string) string {
	fmt.Print(label)
	sc.Scan()
	return sc.Text()
}
