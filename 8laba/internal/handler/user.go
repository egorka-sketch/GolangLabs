package handler

import (
	"8laba/internal/repository"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Birthday string `json:"birth_day"`
}
type Handler struct {
	repo repository.UserRepository
}
type UserRequest struct {
	Name     string `json:"name"`
	Birthday string `json:"birth_day"`
}
type ValErrors struct {
	Msg string `json:"msg"`
}

func (e ValErrors) Error() string {
	return e.Msg
}

func toUserResp(u repository.User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Birthday: u.Birthday.Format("2006-01-02"),
	}
}
func NewHandler(r repository.UserRepository) *Handler {
	return &Handler{
		repo: r,
	}
}
func writeJson(w http.ResponseWriter, status int, d any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(d)
}
func writeJsonError(w http.ResponseWriter, status int, msg string) {
	writeJson(w, status, map[string]string{"error": msg})
}

func (h *Handler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	limit, err := queryInt(r, "limit", 20)
	if err != nil {
		handleError(w, ValErrors{Msg: "limit должен быть числом"})
		return
	}
	if limit < 1 || 100 < limit {
		handleError(w, ValErrors{Msg: "limit должен быть от 1 до 100"})
		return
	}
	offset, err := queryInt(r, "offset", 0)
	if err != nil {
		handleError(w, ValErrors{Msg: "offset должен быть числом"})
		return
	}
	if offset < 0 {
		handleError(w, ValErrors{Msg: "offset не может быть отрицательным"})
		return
	}
	var minAge, maxAge int
	if minAge, err = queryInt(r, "min_age", -1); err != nil || minAge < -1 {
		handleError(w, ValErrors{Msg: "min_age должен быть неотрицательным"})
		return
	}
	if maxAge, err = queryInt(r, "max_age", -1); err != nil || maxAge < -1 {
		handleError(w, ValErrors{Msg: "max_age должен быть неотрицательным"})
		return
	}
	if minAge != -1 && maxAge != -1 && minAge > maxAge {
		handleError(w, ValErrors{Msg: "min_age не может быть больше чем max_age"})
		return
	}
	filter := repository.UserFilter{Limit: limit, Offset: offset, Name: strings.TrimSpace(r.URL.Query().Get("name")), MinAge: minAge, MaxAge: maxAge}
	users, err := h.repo.GetAll(r.Context(), filter)
	if err != nil {
		handleError(w, err)
		return
	}
	resp := make([]UserResponse, len(users))
	for i, u := range users {
		resp[i] = toUserResp(u)
	}
	writeJson(w, http.StatusOK, resp)
}

func (h *Handler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, ValErrors{Msg: "некорректный id"})
		return
	}
	user, err := h.repo.GetByID(r.Context(), userId)

	if err != nil {
		handleError(w, err)
		return
	}
	userResp := toUserResp(user)
	writeJson(w, http.StatusOK, userResp)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		handleError(w, ValErrors{Msg: "некорректный json"})
		return
	}
	u, err := validate(req)
	if err != nil {
		handleError(w, err)
		return
	}
	user, err := h.repo.Create(r.Context(), u)
	if err != nil {
		handleError(w, err)
		return
	}
	userResp := toUserResp(user)
	writeJson(w, http.StatusCreated, userResp)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, ValErrors{Msg: "некорректный id"})
		return
	}
	var req UserRequest
	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		handleError(w, ValErrors{Msg: "некорректный json"})
		return
	}
	u, err := validate(req)
	if err != nil {
		handleError(w, err)
		return
	}
	u.ID = userId
	user, err := h.repo.Update(r.Context(), u)
	if err != nil {
		handleError(w, err)
		return
	}
	userResp := toUserResp(user)
	writeJson(w, http.StatusOK, userResp)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, ValErrors{Msg: "некорректный id"})
		return
	}
	err = h.repo.Delete(r.Context(), userId)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validate(req UserRequest) (repository.User, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return repository.User{}, ValErrors{Msg: "имя не может быть пустым"}
	}
	if utf8.RuneCountInString(name) > 30 {
		return repository.User{}, ValErrors{Msg: "имя не может быть более 30 символов"}
	}
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return repository.User{}, ValErrors{Msg: "дата рождения должна быть в формате ГГГГ-ММ-ДД"}
	}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if birthday.After(today) {
		return repository.User{}, ValErrors{Msg: "дата не может быть в будущем"}
	}
	if birthday.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return repository.User{}, ValErrors{Msg: "дата не может быть раньше начала 1900 года"}
	}
	u := repository.User{Name: name, Birthday: birthday}
	return u, nil
}
func handleError(w http.ResponseWriter, err error) {
	var ve ValErrors
	if errors.As(err, &ve) {
		writeJsonError(w, http.StatusBadRequest, ve.Error())
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		writeJsonError(w, http.StatusNotFound, "пользователь не найден")
		return
	}
	log.Println("внутренняя ошибка", err)
	writeJsonError(w, http.StatusInternalServerError, "внутренняя ошибка")
}
func queryInt(r *http.Request, key string, def int) (int, error) {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def, nil
	}

	return strconv.Atoi(s)
}
