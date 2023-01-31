package catchall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) SetupSuite() {
}

func (suite *TestSuite) TestJWTGet() {
	testData, _ := os.ReadFile("testdata/jwt.json")
	testDataStruct := &JSONStruct{}
	json.Unmarshal([]byte(testData), testDataStruct)
	k8sServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/realms/redhat-external/" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(testData)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer k8sServer.Close()

	os.Setenv("KEYCLOAK_SERVER", k8sServer.URL)

	// dummy muxer for the test
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(MakeNewMBOPServer().MainHandler))

	sut := httptest.NewServer(mux)
	defer sut.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/jwt", sut.URL))
	b, _ := io.ReadAll(resp.Body)

	assert.Nil(suite.T(), err, "error was not nil")
	assert.Equal(suite.T(), 200, resp.StatusCode, "status code not good")
	assert.Equal(suite.T(), testDataStruct.PublicKey, string(b), fmt.Sprintf("expected body doesn't match: %v", string(b)))

	defer resp.Body.Close()
}

func (suite *TestSuite) TestGetUrl() {
	os.Setenv("KEYCLOAK_SERVER", "http://test")
	path := MakeNewMBOPServer().getURL("path", map[string]string{"hi": "you"})
	assert.Equal(suite.T(), "http://test/path?hi=you", path, "did not match")
}

func (suite *TestSuite) TearDownSuite() {
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
