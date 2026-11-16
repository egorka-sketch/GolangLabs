package auth

import (
	"encoding/json"
	"net/http"
	"strings"
)

func BearerToken(r *http.Request) (string, bool) {
	token, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !found || token == "" {
		return "", false
	}
	return token, true
}

func RequireAuth(store *SessionStore, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := BearerToken(r)
		if ok != true {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "требуется авторизация"})
			return
		}
		_, err := store.Get(token)
		if err != true {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "требуется авторизация"})
			return
		}
		next.ServeHTTP(w, r)
	}
}
