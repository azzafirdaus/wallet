package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Middleware -> http middleware
func Middleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		token := r.Header.Get("Authorization")

		sessionID := getSessionByToken(token)
		userID, status := checkSession(sessionID)
		if !status {
			response := Response{
				Status: statusFail,
				Data: ResponseError{
					Error: "Authorization failed",
				},
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)

			return
		}
		r.ParseForm()
		r.Form.Set("user_id", userID)

		next(w, r, ps)
		return
	}
}

func getSessionByToken(token string) (sessionID string) {
	arr := strings.Fields(token)
	if len(arr) != 2 {
		return
	}

	if strings.ToLower(arr[0]) != "token" {
		return
	}

	sessionID = arr[1]

	return
}

func checkSession(sessionID string) (userID string, status bool) {
	userID, err := getSession(database, sessionID)
	if err != nil {
		return
	}

	if userID == "" {
		return
	}

	status = true
	return
}
