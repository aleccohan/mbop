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
	testDataStruct *jwk2pem.JWKeys
	testData       []byte
}

func (suite *TestSuite) SetupTest() {
	testData, _ := os.ReadFile("testdata/jwt.json")
	testDataStruct := &jwk2pem.JWKeys{}
	err := json.Unmarshal([]byte(testData), testDataStruct)
	assert.Nil(suite.T(), err, "error was not nil")

	suite.testDataStruct = testDataStruct
	suite.testData = testData

	os.Setenv("JWT_MODULE", "aws")
}

func (suite *TestSuite) TestAwsJWTGetNoKid() {
	suite.T().Skip()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(suite.testData)
	}))
	defer mockServer.Close()
	os.Setenv("JWK_URL", fmt.Sprintf("%s/v1/jwt", mockServer.URL))

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

func (suite *TestSuite) TestAwsJWTGetNoKidMatch() {
	suite.T().Skip()
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(JWTV1Handler))

	sut := httptest.NewServer(mux)
	defer sut.Close()

	kid := "abc123"
	resp, err := http.Get(fmt.Sprintf("%s/v1/jwt?kid=%s", sut.URL, kid))
	b, _ := io.ReadAll(resp.Body)

	assert.Nil(suite.T(), err, "error was not nil")
	assert.Equal(suite.T(), 400, resp.StatusCode, "status code not good")
	assert.Equal(suite.T(), fmt.Sprintf("no JWK for kid: %s\n", kid), string(b), fmt.Sprintf("expected body doesn't match: %v", string(b)))

	defer resp.Body.Close()
}

func (suite *TestSuite) TestAwsJWTGetPositiveKidMatch() {
	suite.T().Skip()
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(JWTV1Handler))

	sut := httptest.NewServer(mux)
	defer sut.Close()

	kid := "b4OUzJFABPSRwxX5VN7lYswVj9qoc3tet0tsfG5MSME"
	resp, err := http.Get(fmt.Sprintf("%s/v1/jwt?kid=%s", sut.URL, kid))
	b, _ := io.ReadAll(resp.Body)

	assert.Nil(suite.T(), err, "error was not nil")
	assert.Equal(suite.T(), 200, resp.StatusCode, "status code not good")
	assert.Equal(suite.T(), "", string(b), fmt.Sprintf("expected body doesn't match: %v", string(b)))

	defer resp.Body.Close()
}

func (suite *TestSuite) TearDownSuite() {
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
