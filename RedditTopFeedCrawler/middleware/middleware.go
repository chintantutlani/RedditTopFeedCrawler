package middleware

import "net/http"

type Middleware struct{}

func (m *Middleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authenticated := performAuthentication(r)
		if !authenticated {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func performAuthentication(r *http.Request) bool {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	expectedUsername := "admin"
	expectedPassword := "password"

	return username == expectedUsername && password == expectedPassword
}
