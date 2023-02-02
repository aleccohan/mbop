package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/RedHatInsights/jwk2pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) SetupSuite() {
}

func (suite *TestSuite) TestAwsJWTGetNoKid() {
	testData, _ := os.ReadFile("testdata/jwt.json")
	testDataStruct := &jwk2pem.JWKeys{}
	err := json.Unmarshal([]byte(testData), testDataStruct)
	assert.Nil(suite.T(), err, "error was not nil")

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(testData)
	}))
	defer mockServer.Close()

	os.Setenv("JWT_MODULE", "aws")
	os.Setenv("JWK_URL", fmt.Sprintf("%s/v1/jwt", mockServer.URL))

	// dummy muxer for the test
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(JWTV1Handler))

	sut := httptest.NewServer(mux)
	defer sut.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/jwt", sut.URL))
	b, _ := io.ReadAll(resp.Body)

	assert.Nil(suite.T(), err, "error was not nil")
	assert.Equal(suite.T(), 400, resp.StatusCode, "status code not good")
	assert.Equal(suite.T(), "kid required to return correct pub key\n", string(b), fmt.Sprintf("expected body doesn't match: %v", string(b)))

	defer resp.Body.Close()
}

func (suite *TestSuite) TearDownSuite() {
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
