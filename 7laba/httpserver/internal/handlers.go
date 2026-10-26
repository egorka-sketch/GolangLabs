package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type Status struct {
	User     string `json:"user"`
	Online   bool   `json:"online"`
	LastSeen string `json:"lastSeen,omitempty"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
func DataHandler(w http.ResponseWriter, r *http.Request) {
	var person Person
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&person)
	if err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		log.Println(err)
		return
	}
	if person.Name == "" || person.Age == 0 {
		http.Error(w, "Поля имя и возраст не могут быть пустыми", http.StatusBadRequest)
		return
	}
	log.Printf("Имя: %s Возраст: %d", person.Name, person.Age)
	fmt.Fprintf(w, "Успешная отправка")
}
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	user := r.PathValue("user")
	status := Status{User: user}
	if rand.IntN(2) == 0 {
		status.Online = true
	} else {
		lastOnline := time.Now().Add(-time.Duration(rand.IntN(120)) * time.Minute)
		status.LastSeen = lastOnline.Format("15:04")
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		log.Println("не удалось отправить ответ:", err)
	}
}
