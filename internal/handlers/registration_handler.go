package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/redhatinsights/mbop/internal/store"
)

type registationCreateRequest struct {
	Uid *string `json:"uid,omitempty"`
}

type registrationCreateResponse struct {
	Registered string `json:"registered,omitempty"`
	OrgID      string `json:"org_id,omitempty"`
	UID        string `json:"uid,omitempty"`
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		do500(w, "failed to read body bytes: "+err.Error())
		return
	}

	var body registationCreateRequest
	err = json.Unmarshal(b, &body)
	if err != nil {
		do400(w, "failed to unmarshal body: "+err.Error())
		return
	}

	if body.Uid == nil || *body.Uid == "" {
		do400(w, "required parameter [uid] not found in body")
		return
	}

	raw := r.Header.Get("Authorization")
	if raw == "" {
		do400(w, "need bearer token for registration")
		return
	}

	if !strings.HasPrefix(raw, "Bearer ") {
		do400(w, "need bearer token for registration")
		return
	}

	parts := strings.Fields(raw)
	if len(parts) != 2 {
		do400(w, "bearer token in improper format")
		return
	}

	db := store.GetStore()

	// TODO: look up orgid from token?
	token := parts[1]
	orgId := _getOrgIdFromToken(token)

	_, err = db.Find(orgId, *body.Uid)
	if err == nil {
		doError(w, "existing registration found", 409)
		return
	}

	id, err := db.Create(&store.Registration{
		OrgID: orgId,
		UID:   *body.Uid,
	})
	if err != nil {
		do500(w, "failed to create registration: "+err.Error())
		return
	}

	sendJSONWithStatusCode(w, &registrationCreateResponse{
		Registered: id,
		OrgID:      orgId,
		UID:        *body.Uid,
	}, 201)
}

// TODO: call out to AMS
func _getOrgIdFromToken(token string) string {
	return token
}
