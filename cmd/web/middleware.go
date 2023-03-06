package main

import "net/http"

func (a *application) IsAuthenticated(r *http.Request) bool {
	return a.SessionManager.Exists(r.Context(), "userID")
}

func (a *application) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}
