package internal

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

func SendMsg(addr string) error {

	conn, err := net.Dial("tcp", addr)
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
