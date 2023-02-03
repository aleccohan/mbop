package handlers

import (
	"encoding/json"
	"net/http"

	l "github.com/redhatinsights/mbop/internal/logger"
)

func sendJSON(w http.ResponseWriter, data any) {
	sendJSONWithStatusCode(w, data, 200)
}

func sendJSONWithStatusCode(w http.ResponseWriter, data any, code int) {
	b, _ := json.Marshal(data)

	w.WriteHeader(code)
	_, err := w.Write(b)
	if err != nil {
		l.Log.Error(err, "error writing response")
	}
}

func do500(w http.ResponseWriter, msg string) {
	doError(w, msg, 500)
}

func do400(w http.ResponseWriter, msg string) {
	doError(w, msg, 400)
}

func do404(w http.ResponseWriter, msg string) {
	doError(w, msg, 404)
}

func doError(w http.ResponseWriter, msg string, code int) {
	sendJSONWithStatusCode(w, newResponse(msg), code)
}
