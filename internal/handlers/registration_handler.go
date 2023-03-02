package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/redhatinsights/mbop/internal/store"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type registationCreateRequest struct {
	UID *string `json:"uid,omitempty"`
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

	gatewayCN := r.Header.Get("x-rh-certauth-cn")
	parts := strings.Split(gatewayCN, "=")
	if gatewayCN == "" || len(parts) < 2 {
		doError(w, "[x-rh-certauth-cn] header not present", 400)
		return
	}

	if parts[1] != *body.UID {
		doError(w, "x-rh-certauth-cn does not match uid", 400)
		return
	}

	db := store.GetStore()
	_, err = db.Find(id.Identity.OrgID, *body.UID)
	if err == nil {
		doError(w, "existing registration found", 409)
		return
	}

	_, err = db.Create(&store.Registration{
		OrgID: id.Identity.OrgID,
		UID:   *body.UID,
	})
	if err != nil {
		if errors.Is(store.ErrUIDAlreadyExists, err) {
			doError(w, "existing registration found", 409)
		} else {
			do500(w, "failed to create registration: "+err.Error())
		}
		return
	}

	sendJSONWithStatusCode(w, newResponse("Successfully registered"), 201)
}
