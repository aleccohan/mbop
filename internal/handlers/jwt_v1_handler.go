package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/RedHatInsights/jwk2pem"
	"io"
	"net/http"
	"os"
)

func JWTV1Handler(w http.ResponseWriter, r *http.Request) {
	switch os.Getenv("JWT_MODULE") {
	case "aws":
		kid := r.URL.Query().Get("kid")
		if kid == "" {
			http.Error(w, "kid required to return correct pub key", http.StatusBadRequest)
			return
		}

		type JwkRsp struct {
			PublicKey string `json:"public_key"`
		}

		jwkUrl := os.Getenv("JWK_URL")
		resp, err := http.Get(jwkUrl)

		if err != nil {
			http.Error(w, "could not get JWKs", http.StatusBadRequest)
			return
		}

		defer resp.Body.Close()

		bdata, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "could not read response", http.StatusInternalServerError)
			return
		}

		keys := jwk2pem.JWKeys{}
		json.Unmarshal([]byte(bdata), &keys)

		pem := jwk2pem.JWKsToPem(keys, kid)

		if pem == nil {
			http.Error(w, fmt.Sprintf("no JWK for kid: %v", kid), http.StatusInternalServerError)
			return
		}

		w.Write(pem)
	default:
		fmt.Println("CATCHALL")
		CatchAll(w, r)
	}
}
