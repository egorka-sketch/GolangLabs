package internal

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type Center struct {
	mu    sync.Mutex
	users map[*websocket.Conn]struct{}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (c *Center) ChatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("не удалось открыть WebSocket:", err)
		return
	}
	defer conn.Close()
	log.Println("В чате новый пользователь: ", conn.RemoteAddr())
	c.join(conn)
	defer c.leave(conn)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("сбой работы чата пользователя", conn.RemoteAddr())
			} else {
				log.Println("Пользователь отключился", conn.RemoteAddr())
			}
			return
		}
		log.Println("сообщение :", string(msg))
		c.announce(msg)
	}
}

func NewUser() *Center {
	return &Center{
		users: make(map[*websocket.Conn]struct{}),
	}
}
func (c *Center) join(conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.users[conn] = struct{}{}
	log.Println("Пользователей в чате", len(c.users))
}
func (c *Center) leave(conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.users, conn)
	log.Println("Пользователей в чате", len(c.users))
}

func (c *Center) announce(msg []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for conn := range c.users {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("не дошло до", conn.RemoteAddr(), err)
			conn.Close()
		}
	}
}
func (c *Center) CloseAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	bye := websocket.FormatCloseMessage(websocket.CloseGoingAway, "сервер закрывается")
	for conn := range c.users {
		if err := conn.WriteControl(websocket.CloseMessage, bye, time.Now().Add(time.Second)); err != nil {
			log.Println("не удалось отключить пользователя", conn.RemoteAddr(), err)
		}
		conn.Close()
	}
}
