package handler

import (
	"9laba/internal/server/auth"
	"9laba/internal/server/repository"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"unicode/utf8"
)

type UserData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type AuthHandler struct {
	accounts repository.AccountRepo
	sessions *auth.SessionStore
}
type Response struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

func NewAuthHandler(accounts repository.AccountRepo, sessions *auth.SessionStore) *AuthHandler {
	return &AuthHandler{
		accounts: accounts,
		sessions: sessions,
	}
}
func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var user UserData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, "некорректный json")
		return
	}
	login := strings.TrimSpace(user.Login)
	password := user.Password
	if utf8.RuneCountInString(login) < 6 {
		writeJsonError(w, http.StatusBadRequest, "слишком короткий логин, минимальная длина 6 символов")
		return
	}
	if utf8.RuneCountInString(password) < 6 {
		writeJsonError(w, http.StatusBadRequest, "слишком короткий пароль, минимальная длина 6 символов")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}
	acc, err := h.accounts.Create(r.Context(), login, string(hash))
	if errors.Is(err, repository.ErrLoginTaken) {
		writeJsonError(w, http.StatusConflict, "логин занят")
		return
	}
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	writeJson(w, http.StatusCreated, Response{ID: acc.ID, Login: acc.Login, Role: acc.Role})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user UserData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, "некорректный json")
		return
	}
	login := strings.TrimSpace(user.Login)
	acc, err := h.accounts.Get(r.Context(), login)
	if errors.Is(err, repository.ErrNotFound) {
		writeJsonError(w, http.StatusUnauthorized, "неверный логин или пароль")
		return
	}
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(user.Password))
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "неверный логин или пароль")
		return
	}
	token, err := h.sessions.Create(acc.ID, acc.Role)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}
	writeJson(w, http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token, _ := auth.BearerToken(r)
	h.sessions.Delete(token)
	w.WriteHeader(http.StatusNoContent)
}
