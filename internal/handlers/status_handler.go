package handlers

import (
	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/models"
	"net/http"
)

func Status(w http.ResponseWriter, r *http.Request) {
	status := models.Status{
		ConfiguredModules: models.ConfiguredModules{
			Users:  config.Get().UsersModule,
			Mailer: config.Get().MailerModule,
			JWT:    config.Get().JwtModule,
		},
	}

	sendJSON(w, status)
}
