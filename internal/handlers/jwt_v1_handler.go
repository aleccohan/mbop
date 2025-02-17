package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/RedHatInsights/jwk2pem"
	l "github.com/redhatinsights/mbop/internal/logger"
)

type JWTResp struct {
	Pubkey string `json:"pubkey"`
}

func JWTV1Handler(w http.ResponseWriter, r *http.Request) {
	switch os.Getenv("JWT_MODULE") {
	case "aws":
		kid := r.URL.Query().Get("kid")
		if kid == "" {
			do400(w, "kid required to return correct pub key")
			return
		}

		jwkURL := os.Getenv("JWK_URL")
		resp, err := http.Get(jwkURL) //nolint

		if err != nil {
			l.Log.Error(err, "error getting JWKs")
			do500(w, "error getting JWKs: "+err.Error())
			return
		}

		defer resp.Body.Close()

		bdata, err := io.ReadAll(resp.Body)
		if err != nil {
			l.Log.Error(err, "error reading JWKs")
			do500(w, "error reading JWKs: "+err.Error())
			return
		}

		keys := jwk2pem.JWKeys{}
		err = json.Unmarshal([]byte(bdata), &keys)
		if err != nil {
			do400(w, "failed to parse response: "+err.Error())
			return
		}

		pem := jwk2pem.JWKsToPem(keys, kid)

		if pem == nil {
			do404(w, "no JWK for kid: "+kid)
			return
		}

		sendJSON(w, JWTResp{Pubkey: strings.TrimSuffix(string(pem), "\n")})
	default:
		CatchAll(w, r)
	}
}
