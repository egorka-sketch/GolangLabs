package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
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
