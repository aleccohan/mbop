package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/redhatinsights/mbop/internal/store"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type registationCreateRequest struct {
	UID *string `json:"uid,omitempty"`
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

	if body.UID == nil || *body.UID == "" {
		do400(w, "required parameter [uid] not found in body")
		return
	}

	id := identity.Get(r.Context())
	if !id.Identity.User.OrgAdmin {
		doError(w, "user must be org admin to register satellite", 403)
		return
	}

	db := store.GetStore()
	_, err = db.Find(id.Identity.OrgID, *body.UID)
	if err == nil {
		doError(w, "existing registration found", 409)
		return
	}

	guid, err := db.Create(&store.Registration{
		OrgID: id.Identity.OrgID,
		UID:   *body.UID,
	})
	if err != nil {
		do500(w, "failed to create registration: "+err.Error())
		return
	}

	sendJSONWithStatusCode(w, &registrationCreateResponse{
		Registered: guid,
		OrgID:      id.Identity.OrgID,
		UID:        *body.UID,
	}, 201)
}
