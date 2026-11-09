package handler

import (
	"8laba/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type fakeRepo struct {
	users []repository.User
	user  repository.User
	err   error
}

var _ repository.UserRepository = (*fakeRepo)(nil)

func (f *fakeRepo) GetByID(ctx context.Context, id int) (repository.User, error) {
	if f.err != nil {
		return repository.User{}, f.err
	}
	return f.user, nil
}
func (f *fakeRepo) GetAll(ctx context.Context, filter repository.UserFilter) ([]repository.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.users, nil
}
func (f *fakeRepo) Create(ctx context.Context, user repository.User) (repository.User, error) {
	if f.err != nil {
		return repository.User{}, f.err
	}
	user.ID = 1
	return user, nil
}
func (f *fakeRepo) Update(ctx context.Context, user repository.User) (repository.User, error) {
	if f.err != nil {
		return repository.User{}, f.err
	}
	return user, nil
}
func (f *fakeRepo) Delete(ctx context.Context, id int) error {
	return f.err
}

var ivan = repository.User{ID: 1, Name: "Иван", Birthday: time.Date(2001, 3, 15, 0, 0, 0, 0, time.UTC)}
var dbDown = errors.New("connection refused")

const validBody = `{"name":"Иван","birth_day":"2001-03-15"}`

func TestGetByID(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		repo       *fakeRepo
		wantStatus int
	}{
		{"найден", "/users/1", &fakeRepo{user: ivan}, http.StatusOK},
		{
			"кривой id", "/users/abc", &fakeRepo{}, http.StatusBadRequest,
		},
		{
			"не найден", "/users/999", &fakeRepo{err: repository.ErrNotFound}, http.StatusNotFound,
		},
		{
			"база упала", "/users/1", &fakeRepo{err: dbDown}, http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.repo)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /users/{id}", h.HandleGetById)
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("статус = %d, ожидали %d; тело: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantStatus == http.StatusOK {
				var got UserResponse
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("ответ не JSON: %v", err)
				}
				want := UserResponse{ID: 1, Name: "Иван", Birthday: "2001-03-15"}
				if got != want {
					t.Errorf("тело = %+v, ожидали %+v", got, want)
				}
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		repo       *fakeRepo
		wantStatus int
	}{
		{"без параметров", "/users", &fakeRepo{users: []repository.User{ivan}}, http.StatusOK},
		{"пустой список", "/users", &fakeRepo{}, http.StatusOK},
		{"все фильтры сразу", "/users?name=iv&min_age=20&max_age=30&limit=10&offset=0", &fakeRepo{}, http.StatusOK},
		{"только min_age", "/users?min_age=20", &fakeRepo{}, http.StatusOK},
		{"limit не число", "/users?limit=abc", &fakeRepo{}, http.StatusBadRequest},
		{"limit = 0", "/users?limit=0", &fakeRepo{}, http.StatusBadRequest},
		{"limit = 101", "/users?limit=101", &fakeRepo{}, http.StatusBadRequest},
		{"отрицательный offset", "/users?offset=-1", &fakeRepo{}, http.StatusBadRequest},
		{"min_age не число", "/users?min_age=abc", &fakeRepo{}, http.StatusBadRequest},
		{"отрицательный min_age", "/users?min_age=-5", &fakeRepo{}, http.StatusBadRequest},
		{"отрицательный max_age", "/users?max_age=-5", &fakeRepo{}, http.StatusBadRequest},
		{"min_age больше max_age", "/users?min_age=30&max_age=20", &fakeRepo{}, http.StatusBadRequest},
		{"база упала", "/users", &fakeRepo{err: dbDown}, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.repo)
			mux := http.NewServeMux()
			mux.HandleFunc("GET /users", h.HandleGetAll)
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("статус = %d, ожидали %d; тело: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantStatus == http.StatusOK {
				var got []UserResponse
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("ответ не JSON: %v", err)
				}
				if len(got) != len(tt.repo.users) {
					t.Fatalf("%d пользователей. ожидали %d", len(got), len(tt.repo.users))
				}
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repo       *fakeRepo
		wantStatus int
	}{
		{"создан", validBody, &fakeRepo{}, http.StatusCreated},
		{"кривой json", `{"name":`, &fakeRepo{}, http.StatusBadRequest},
		{"пустое имя", `{"name":"","birth_day":"2001-03-15"}`, &fakeRepo{}, http.StatusBadRequest},
		{"имя из пробелов", `{"name":"   ","birth_day":"2001-03-15"}`, &fakeRepo{}, http.StatusBadRequest},
		{"имя ровно 30 символов", `{"name":"` + strings.Repeat("я", 30) + `","birth_day":"2001-03-15"}`, &fakeRepo{}, http.StatusCreated},
		{"имя из 31 символа", `{"name":"` + strings.Repeat("я", 31) + `","birth_day":"2001-03-15"}`, &fakeRepo{}, http.StatusBadRequest},
		{"неверный формат даты", `{"name":"Иван","birth_day":"15.03.2001"}`, &fakeRepo{}, http.StatusBadRequest},
		{"дата в будущем", `{"name":"Иван","birth_day":"2999-01-01"}`, &fakeRepo{}, http.StatusBadRequest},
		{"ровно 1900-01-01", `{"name":"Иван","birth_day":"1900-01-01"}`, &fakeRepo{}, http.StatusCreated},
		{"1899-12-31", `{"name":"Иван","birth_day":"1899-12-31"}`, &fakeRepo{}, http.StatusBadRequest},
		{"база упала", validBody, &fakeRepo{err: dbDown}, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.repo)
			mux := http.NewServeMux()
			mux.HandleFunc("POST /users", h.CreateUser)
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("статус = %d, ожидали %d; тело: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantStatus == http.StatusCreated {
				var got UserResponse
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("ответ не JSON: %v", err)
				}
				if got.ID != 1 {
					t.Fatalf("ID %d. ожидали 1", got.ID)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		body       string
		repo       *fakeRepo
		wantStatus int
	}{
		{"обновлён", "/users/5", validBody, &fakeRepo{}, http.StatusOK},
		{"кривой id", "/users/abc", validBody, &fakeRepo{}, http.StatusBadRequest},
		{"кривой json", "/users/5", `{"name":`, &fakeRepo{}, http.StatusBadRequest},
		{"пустое имя", "/users/5", `{"name":"","birth_day":"2001-03-15"}`, &fakeRepo{}, http.StatusBadRequest},
		{"не найден", "/users/999", validBody, &fakeRepo{err: repository.ErrNotFound}, http.StatusNotFound},
		{"база упала", "/users/5", validBody, &fakeRepo{err: dbDown}, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.repo)
			mux := http.NewServeMux()
			mux.HandleFunc("PUT /users/{id}", h.UpdateUser)
			req := httptest.NewRequest(http.MethodPut, tt.url, strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("статус = %d, ожидали %d; тело: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantStatus == http.StatusOK {
				var got UserResponse
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("ответ не JSON: %v", err)
				}
				if got.ID != 5 {
					t.Errorf("id = %d, ожидали 5", got.ID)
				}
			}
		})

	}
}
func TestDelete(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		repo       *fakeRepo
		wantStatus int
	}{
		{"удалён", "/users/1", &fakeRepo{}, http.StatusNoContent},
		{"кривой id", "/users/abc", &fakeRepo{}, http.StatusBadRequest},
		{"не найден", "/users/999", &fakeRepo{err: repository.ErrNotFound}, http.StatusNotFound},
		{"база упала", "/users/1", &fakeRepo{err: dbDown}, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.repo)
			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /users/{id}", h.DeleteUser)
			req := httptest.NewRequest(http.MethodDelete, tt.url, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("статус = %d, ожидали %d; тело: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}

		})
	}
}
