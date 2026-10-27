package internal

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("не удалось открыть WebSocket:", err)
		return
	}
	defer conn.Close()
	log.Println("В чате новый пользователь: ", conn.RemoteAddr())

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("ошибка чтения сообщения:", err)
			return
		}
		log.Println("сообщение :", string(msg))
		if err = conn.WriteMessage(mt, msg); err != nil {
			log.Println("сообщение не отправлено: ", err)
			return
		}
	}
}
