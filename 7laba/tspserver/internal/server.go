package internal

import (
	"log"
	"net"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}
func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	log.Println("Слушаю порт", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Подключился", conn.RemoteAddr())
		handle(conn)
	}
}
