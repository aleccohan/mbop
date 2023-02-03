package handlers

import (
	"net/http"
	"os"

	"github.com/redhatinsights/mbop/internal/models"
)

func Status(w http.ResponseWriter, r *http.Request) {
	status := models.Status{
		ConfiguredModules: models.ConfiguredModules{
			Users:  os.Getenv("USERS_MODULE"),
			Mailer: os.Getenv("MAILER_MODULE"),
			JWT:    os.Getenv("JWT_MODULE"),
		},
	}

	_, _ = w.Write(status.ToJSON())
}
