package auth

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const CookieName = "token"

func TokenFromRequest(r *http.Request) (string, bool) {
	if token, ok := BearerToken(r); ok {
		return token, true
	}
	if c, err := r.Cookie(CookieName); err == nil && c.Value != "" {
		return c.Value, true
	}
	return "", false
}

func BearerToken(r *http.Request) (string, bool) {
	token, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !found || token == "" {
		return "", false
	}
	return token, true
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func RequireAuth(jm *JWTManager, store *SessionStore, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := TokenFromRequest(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "требуется авторизация")
			return
		}
		claims, err := jm.Parse(token)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "недействительный токен")
			return
		}
		if _, ok := store.Get(claims.ID); !ok {
			writeJSONError(w, http.StatusUnauthorized, "сессия завершена")
			return
		}
		accountID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "недействительный токен")
			return
		}
		p := Principal{AccountID: accountID, Role: claims.Role, SessionID: claims.ID}
		next(w, r.WithContext(WithPrincipal(r.Context(), p)))
	}
}
func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := PrincipalFrom(r.Context())
		if !ok || p.Role != role {
			writeJSONError(w, http.StatusForbidden, "недостаточно прав")
			return
		}
		next(w, r)

	}

}
