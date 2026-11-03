package handler

import (
	"8laba/internal/repository"
	"encoding/json"
	"log"
	"net/http"
)

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Birthday string `json:"birth_day"`
}
type Handler struct {
	repo repository.UserRepository
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
	users, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Println("ошибка вывода", err)
		writeJsonError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	resp := make([]UserResponse, len(users))
	for i, u := range users {
		resp[i] = toUserResp(u)
	}
	writeJson(w, http.StatusOK, resp)
}
