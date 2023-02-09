package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"os"

	"github.com/RedHatInsights/jwk2pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	testData       []byte
	testDataStruct *jwk2pem.JWKeys
	mockServer     *httptest.Server
	testPem        []byte
}

func (suite *TestSuite) SetupSuite() {
	suite.testData, _ = os.ReadFile("testdata/jwks.json")
	suite.testDataStruct = &jwk2pem.JWKeys{}
	err := json.Unmarshal([]byte(suite.testData), suite.testDataStruct)
	assert.Nil(suite.T(), err, "error was not nil")
	suite.testPem, _ = os.ReadFile("testdata/pem.json")
}
