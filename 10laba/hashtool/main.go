package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"fmt"
	"hash"
	"os"
	"strings"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n1. Вычислить хеш\n2. Проверить целостность\nВыбор: ")
		if !sc.Scan() {
			return
		}
		choice := sc.Text()
		if choice != "1" && choice != "2" {
			fmt.Println("Нет такого пункта")
			return
		}

		fmt.Print("Текст: ")
		sc.Scan()
		text := sc.Text()
		fmt.Print("Алгоритм: ")
		sc.Scan()
		alg := sc.Text()

		if choice == "1" {
			sum, err := computeHash(alg, text)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(sum)
		} else {
			fmt.Print("Ожидаемый хеш: ")
			sc.Scan()
			expected := sc.Text()
			ok, err := EqualHash(alg, text, expected)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if ok {
				fmt.Println("Совпадает")
			} else {
				fmt.Println("Не совпадает")
			}
		}
	}
}

func computeHash(alg, text string) (string, error) {
	algs := map[string]func() hash.Hash{
		"md5":      md5.New,
		"sha256":   sha256.New,
		"sha512":   sha512.New,
		"sha3-256": func() hash.Hash { return sha3.New256() },
	}
	newHash, ok := algs[alg]
	if !ok {
		return "", fmt.Errorf("нет такого алгоритма: %s", alg)
	}
	h := newHash()
	h.Write([]byte(text))
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}

func EqualHash(alg, text, expected string) (bool, error) {
	got, err := computeHash(alg, text)
	if err != nil {
		return false, err
	}
	expected = strings.ToLower(strings.TrimSpace(expected))
	return hmac.Equal([]byte(got), []byte(expected)), nil
}
