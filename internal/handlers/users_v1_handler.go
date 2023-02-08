package handlers

import (
	"net/http"

	"github.com/redhatinsights/mbop/internal/config"
)

func UsersV1Handler(w http.ResponseWriter, r *http.Request) {
	switch config.Get().UsersModule {
	case awsModule:
		// acctNum := chi.URLParam(r, "accountNumber")
		// users, err := users.ListUsersV1(acctNum)
		// ...
		// w.Write(users.ToJSON())
	default:
		// mbop server instance injected somewhere
		// pass right through to the current handler
		CatchAll(w, r)
	}
}
