package handlers

import (
	"context"
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
		q := initUserQuery(r)

		if !stringInSlice(q.SortOrder, validSortOrder) {
			do400(w, "sortOrder must be one of "+strings.Join(validSortOrder, ", "))
		}

		if !stringInSlice(q.QueryBy, validQueryBy) {
			do400(w, "queryBy must be one of "+strings.Join(validQueryBy, ", "))
		}

		connection, err := sdk.InitSdkConnection(ctx)

		if err != nil {
			do500(w, "Can't build connection: "+err.Error())
		}

		// Get list of Accounts
		collection := connection.AccountsMgmt().V1().Accounts().List()

		collection = addQueryOrder(collection, q)

		accountsGetResponse, err := collection.Send()
		if err != nil {
			do500(w, "Cant Retrieve Accounts: "+err.Error())
		}

		users := models.Users{}
		if accountsGetResponse.Items().Empty() {
			sendJSON(w, users)
		}

		users, err = usersV1.ResponseToUsers(accountsGetResponse, connection)

		if err != nil {
			do500(w, "Cant Retrieve Role Bindings: "+err.Error())
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
