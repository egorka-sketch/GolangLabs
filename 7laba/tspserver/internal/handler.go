package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func handle(conn net.Conn) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		msg := sc.Text()
		log.Printf("Сообщение %s от %v", msg, conn.RemoteAddr())
		fmt.Fprintln(conn, "Сообщение принято:", msg)
	}
	if err := sc.Err(); err != nil {
		log.Printf("%v Сбой при приеме сообщения %v", conn.RemoteAddr(), err)
		return
	}
	log.Printf("Клиент %v отключился", conn.RemoteAddr())
}
