package middleware

import "net/http"

type contextKey string

const requestIDKey contextKey = "request_id"

func GetRequestID(r *http.Request) string {
	value := r.Context().Value(requestIDKey)
	id, ok := value.(string)
	if !ok {
		return ""
	}
	return id
}

