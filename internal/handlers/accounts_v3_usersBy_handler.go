package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/models"
	"github.com/redhatinsights/mbop/internal/service/ocm"
)

func AccountsV3UsersByHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch config.Get().UsersModule {
	case amsModule, mockModule:
		orgID := chi.URLParam(r, "orgID")
		if orgID == "" {
			do400(w, "Request URL must include orgID: /v3/accounts/{orgID}/usersBy")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			do500(w, "failed to read request body: "+err.Error())
			return
		}
		defer r.Body.Close()

		var usersByBody models.UsersByBody
		err = json.Unmarshal(body, &usersByBody)
		if err != nil {
			do400(w, "failed to parse request body: "+err.Error()+", request must include 'primaryEmail', 'emailStartsWith', or 'principalStartsWith'")
			return
		}

		q, err := initAccountV3UserQuery(r)
		if err != nil {
			do400(w, err.Error())
			return
		}

		// Create new SDK client
		client, err := ocm.NewOcmClient()
		if err != nil {
			do400(w, err.Error())
			return
		}

		err = client.InitSdkConnection(ctx)
		if err != nil {
			do500(w, "Can't build connection: "+err.Error())
			return
		}

		u, err := client.GetAccountV3UsersBy(orgID, q, usersByBody)
		if err != nil {
			do500(w, "Cant Retrieve Users: "+err.Error())
			return
		}

		// For each user see if it's an org_admin
		isOrgAdmin, err := client.GetOrgAdmin(u.Users)
		if err != nil {
			do500(w, "Cant Retrieve Role Bindings: "+err.Error())
			return
		}

		for i, user := range u.Users {
			response, ok := isOrgAdmin[user.ID]
			if ok {
				u.Users[i].IsOrgAdmin = response.IsOrgAdmin
			} else {
				user.IsOrgAdmin = false
				if q.AdminOnly { // if admin_only return only the org_admins
					u.RemoveUser(i)
				}
			}
		}

		r := usersToV3Response(u.Users)

		sendJSON(w, r.Responses)
	default:
		// mbop server instance injected somewhere
		// pass right through to the current handler
		CatchAll(w, r)
	}
}
