package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListUsers_QueryAndDecode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/users" {
			t.Errorf("запрос %s %s, ждали GET /users", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("name"); got != "Иван Петров" {
			t.Errorf("name = %q", got)
		}
		if got := r.URL.Query().Get("min_age"); got != "0" {
			t.Errorf("min_age = %q, ждали 0", got)
		}
		if r.URL.Query().Has("max_age") {
			t.Error("max_age не задан, но ушёл в запрос")
		}
		json.NewEncoder(w).Encode([]User{{ID: 1, Name: "Иван Петров", BirthDay: "2001-03-15"}})
	}))
	defer srv.Close()

	zero := 0
	users, err := New(srv.URL).ListUsers(context.Background(), ListFilter{Name: "Иван Петров", MinAge: &zero})
	if err != nil {
		t.Fatalf("ошибка: %v", err)
	}
	if len(users) != 1 || users[0].BirthDay != "2001-03-15" {
		t.Errorf("получили %+v", users)
	}
}

func TestCreateUser_SendsJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q", r.Header.Get("Content-Type"))
		}
		var in UserInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			t.Fatalf("тело не JSON: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(User{ID: 7, Name: in.Name, BirthDay: in.BirthDay})
	}))
	defer srv.Close()

	u, err := New(srv.URL).CreateUser(context.Background(), UserInput{Name: "Иван", BirthDay: "2001-03-15"})
	if err != nil {
		t.Fatalf("ошибка: %v", err)
	}
	if u.ID != 7 || u.Name != "Иван" {
		t.Errorf("получили %+v", u)
	}
}

func TestDeleteUser_NoContent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/users/5" {
			t.Errorf("запрос %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	if err := New(srv.URL).DeleteUser(context.Background(), 5); err != nil {
		t.Fatalf("ошибка: %v", err)
	}
}

func TestAPIErrors(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantMsg string
	}{
		{"404 json", http.StatusNotFound, `{"error":"пользователь не найден"}`, "пользователь не найден"},
		{"401 json", http.StatusUnauthorized, `{"error":"требуется авторизация"}`, "требуется авторизация"},
		{"405 текст", http.StatusMethodNotAllowed, "Method Not Allowed\n", "Method Not Allowed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			_, err := New(srv.URL).GetUser(context.Background(), 1)
			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("ждали *APIError, получили %v", err)
			}
			if apiErr.StatusCode != tt.status || apiErr.Message != tt.wantMsg {
				t.Errorf("получили %d %q", apiErr.StatusCode, apiErr.Message)
			}
		})
	}
}

func TestServerDown(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close()

	_, err := New(url).GetUser(context.Background(), 1)
	if err == nil {
		t.Fatal("ждали ошибку")
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		t.Fatalf("сетевая ошибка не должна быть APIError: %v", err)
	}
}

func TestLogin_TokenInHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/login":
			if r.Header.Get("Authorization") != "" {
				t.Error("до входа заголовка Authorization быть не должно")
			}
			json.NewEncoder(w).Encode(map[string]string{"token": "abc123"})
		case "/users/1":
			if got := r.Header.Get("Authorization"); got != "Bearer abc123" {
				t.Errorf("Authorization = %q", got)
			}
			json.NewEncoder(w).Encode(User{ID: 1})
		case "/auth/logout":
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer srv.Close()

	c := New(srv.URL)
	ctx := context.Background()
	if err := c.Login(ctx, "ivanov", "secret1"); err != nil {
		t.Fatal(err)
	}
	if !c.LoggedIn() {
		t.Fatal("после Login клиент должен считаться вошедшим")
	}
	if _, err := c.GetUser(ctx, 1); err != nil {
		t.Fatal(err)
	}
	if err := c.Logout(ctx); err != nil {
		t.Fatal(err)
	}
	if c.LoggedIn() {
		t.Error("после Logout токен должен стереться")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"неверный логин или пароль"}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	err := c.Login(context.Background(), "ivanov", "wrong")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("ждали 401, получили %v", err)
	}
	if c.LoggedIn() {
		t.Error("при неудачном входе токен сохраняться не должен")
	}
}
