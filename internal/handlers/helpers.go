package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	l "github.com/redhatinsights/mbop/internal/logger"
	"github.com/redhatinsights/mbop/internal/models"
)

var (
	validSortOrder = []string{"asc", "des"}
	validQueryBy   = []string{"userId", "orgId"} // Originally orgId was "principal" but in FedRAMP cluster we only have orgId
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func initUserQuery(r *http.Request) (models.UserQuery, error) {
	q := models.UserQuery{}

	if r.URL.Query().Get("sortOrder") == "" || stringInSlice(r.URL.Query().Get("sortOrder"), validSortOrder) {
		if r.URL.Query().Get("sortOrder") == validSortOrder[1] {
			q.SortOrder = "desc"
		} else {
			q.SortOrder = r.URL.Query().Get("sortOrder")
		}
	} else {
		return q, fmt.Errorf("sortOrder must be one of '', " + strings.Join(validSortOrder, ", "))
	}

	if r.URL.Query().Get("queryBy") == "" || stringInSlice(r.URL.Query().Get("queryBy"), validQueryBy) {
		// Translate bop parameters into AMS parameters
		if r.URL.Query().Get("queryBy") == validQueryBy[0] {
			q.QueryBy = "id"
		}

		if r.URL.Query().Get("queryBy") == validQueryBy[1] {
			q.QueryBy = "organizationId"
		}
	} else {
		return q, fmt.Errorf("queryBy must be one of " + strings.Join(validQueryBy, ", "))
	}

	return q, nil
}
