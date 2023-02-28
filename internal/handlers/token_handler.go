package handlers

import (
	"net/http"
	"time"

	"github.com/redhatinsights/mbop/internal/config"
	l "github.com/redhatinsights/mbop/internal/logger"
	"github.com/redhatinsights/mbop/internal/models"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type TokenResp struct {
	Token string `json:"token"`
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	xrhid := identity.Get(r.Context()).Identity
	if xrhid.OrgID == "" {
		do400(w, "Missing org_id in x-rh-identity")
		return
	}

	if xrhid.User.Username == "" {
		do400(w, "Missing username in x-rh-identity")
		return
	}

	c := config.Get()
	privateKey := []byte(c.PrivateKey)
	pubKey := []byte(c.PublicKey)

	token := models.Token{PrivateKey: privateKey, PublicKey: pubKey}
	ttl, err := time.ParseDuration(c.TokenTTL)
	if err != nil {
		l.Log.Error(err, "Error setting TTL")
	}

	newToken, err := token.Create(ttl, xrhid)
	if err != nil {
		l.Log.Error(err, "Error creating token")
	}

	sendJSON(w, TokenResp{Token: newToken})
}
