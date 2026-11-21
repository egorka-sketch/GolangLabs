package internal

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

func handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			log.Println("Сервер завершает работу")
			fmt.Fprintln(conn, "Сервер завершает работу")
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		case <-done:
		}
	}()
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		msg := sc.Text()
		log.Printf("Сообщение %s от %v", msg, conn.RemoteAddr())
		fmt.Fprintln(conn, "Сообщение принято:", msg)
		log.Printf("ответ отправлен %v", conn.RemoteAddr())
	}
	var ne net.Error
	if err := sc.Err(); err == nil {
		log.Printf("Клиент %v отключился", conn.RemoteAddr())
		return
	} else if errors.As(err, &ne) && ne.Timeout() {
		log.Printf("Клиент %v отключен по истечению времени ожидания", conn.RemoteAddr())
	} else {
		log.Printf("%v Сбой при приеме сообщения %v", conn.RemoteAddr(), err)
	}
}
