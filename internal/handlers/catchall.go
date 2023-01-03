package handlers

import (
	"net/http"

	"github.com/redhatinsights/mbop/internal/service/catchall"
)

// instantiate on startup - only if needed.
var mbop = catchall.MakeNewMBOPServer()

func CatchAll(w http.ResponseWriter, r *http.Request) {
	mbop.MainHandler(w, r)
}
