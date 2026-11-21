package internal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"net"
	"os"
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
	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		return err
	}
	caPEM, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		return err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caPEM)
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13}
	ln, err := tls.Listen("tcp", s.addr, cfg)
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
