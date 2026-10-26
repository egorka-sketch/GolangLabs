package internal

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}
func (s *Server) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	defer ln.Close()

	log.Println("Слушаю порт", s.addr)

	go func() {
		<-ctx.Done()
		ln.Close()
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			log.Println(err)
			continue
		}
		log.Println("Подключился", conn.RemoteAddr())
		wg.Go(func() {
			handle(ctx, conn)
		})
	}
	wg.Wait()
	log.Println("Сервер остановлен")
	return nil
}
