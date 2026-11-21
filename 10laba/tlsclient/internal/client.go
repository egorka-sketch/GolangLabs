package internal

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"
)

func SendMsg(addr string) error {
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		return err
	}
	caPEM, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		return err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) {
		return errors.New("не удалось прочитать ca.crt")
	}
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		ServerName:   "localhost",
		MinVersion:   tls.VersionTLS13}
	conn, err := tls.Dial("tcp", addr, cfg)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("%v Подключился к чату с %v", conn.LocalAddr(), conn.RemoteAddr())
	msg, err := readMsg()
	if err != nil {
		return err
	}
	log.Printf("введено сообщение: %v", msg)
	fmt.Fprintln(conn, msg)
	answer, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}
	log.Printf("ответ сервера: %v", answer)
	return nil

}
func readMsg() (string, error) {
	fmt.Print("Введите сообщение для отправки: ")
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		if err := sc.Err(); err != nil {
			return "", err
		}
		return "", errors.New("пользователь ничего не ввёл")
	}
	return sc.Text(), nil

}
