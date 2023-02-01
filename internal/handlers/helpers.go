package handlers

import (
	"net/http"

	l "github.com/redhatinsights/mbop/internal/logger"
)

func do500(w http.ResponseWriter, msg string) {
	doError(w, msg, 500)
}

func do400(w http.ResponseWriter, msg string) {
	doError(w, msg, 400)
}

func doError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	_, err := w.Write(newResponse(msg))
	if err != nil {
		l.Log.Error(err, "error writing response")
	}
}
