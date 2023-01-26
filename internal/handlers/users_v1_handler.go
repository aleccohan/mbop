package handlers

import (
	"net/http"
	"os"
)

func UsersV1Handler(w http.ResponseWriter, r *http.Request) {
	switch os.Getenv("USERS_MODULE") {
	case "aws":
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
