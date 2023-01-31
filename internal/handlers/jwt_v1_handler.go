package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/RedHatInsights/jwk2pem"
)

func JWTV1Handler(w http.ResponseWriter, r *http.Request) {
	switch os.Getenv("JWT_MODULE") {
	case "aws":
		kid := r.URL.Query().Get("kid")
		if kid == "" {
			http.Error(w, "kid required to return correct pub key", http.StatusBadRequest)
			return
		}

		jwkURL := os.Getenv("JWK_URL")
		resp, err := http.Get(jwkURL) //nolint

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
		err = json.Unmarshal([]byte(bdata), &keys)
		if err != nil {
			http.Error(w, "could not read response", http.StatusInternalServerError)
			return
		}

		pem := jwk2pem.JWKsToPem(keys, kid)

		if pem == nil {
			http.Error(w, fmt.Sprintf("no JWK for kid: %v", kid), http.StatusBadRequest)
			return
		}

		_, err = w.Write(pem)
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Println("CATCHALL")
		CatchAll(w, r)
	}
}
