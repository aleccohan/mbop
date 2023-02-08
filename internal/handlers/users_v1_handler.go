package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/redhatinsights/mbop/internal/config"

	"strings"

	"github.com/redhatinsights/mbop/internal/models"
	usersV1 "github.com/redhatinsights/mbop/internal/service/users_v1"
)

func UsersV1Handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	switch config.Get().UsersModule {
	case awsModule:
		sdk := new(usersV1.OcmSDK)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			do500(w, "failed to read request body: "+err.Error())
			return
		}
		defer r.Body.Close()

		var usernames models.UserBody
		err = json.Unmarshal(body, &usernames)
		if err != nil {
			do400(w, "failed to parse request body: "+err.Error())
			return
		}

		q := initUserQuery(r)
		if !stringInSlice(q.SortOrder, validSortOrder) {
			do400(w, "sortOrder must be one of "+strings.Join(validSortOrder, ", "))
			return
		}

		if !stringInSlice(q.QueryBy, validQueryBy) {
			do400(w, "queryBy must be one of "+strings.Join(validQueryBy, ", "))
			return
		}

		connection, err := sdk.InitSdkConnection(ctx)

		if err != nil {
			do500(w, "Can't build connection: "+err.Error())
			return
		}

		// Get list of Accounts
		search := usersV1.CreateSearchString(usernames)
		collection := connection.AccountsMgmt().V1().Accounts().List().Search(search)

		collection = addQueryOrder(collection, q)
		// Add place here to add usernames into search

		accountsGetResponse, err := collection.Send()
		if err != nil {
			do500(w, "Cant Retrieve Accounts: "+err.Error())
			return
		}

		users := models.Users{}
		if accountsGetResponse.Items().Empty() {
			sendJSON(w, users)
		}

		users, err = usersV1.ResponseToUsers(accountsGetResponse, connection)

		if err != nil {
			do500(w, "Cant Retrieve Role Bindings: "+err.Error())
			return
		}

		// Close SDK Connection
		sdk.CloseSdkConnection(connection)

		sendJSON(w, users)
	default:
		// mbop server instance injected somewhere
		// pass right through to the current handler
		CatchAll(w, r)
	}
}
