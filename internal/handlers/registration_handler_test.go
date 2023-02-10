package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/store"
	"github.com/stretchr/testify/suite"
)

type RegistrationTestSuite struct {
	suite.Suite
	rec   *httptest.ResponseRecorder
	store store.Store
}

func (suite *RegistrationTestSuite) SetupSuite() {
	config.Reset()
	os.Setenv("STORE_BACKEND", "memory")
}

func (suite *RegistrationTestSuite) BeforeTest(_, _ string) {
	suite.rec = httptest.NewRecorder()
	suite.Nil(store.SetupStore())

	suite.store = store.GetStore()
	store.GetStore = func() store.Store { return suite.store }
}

func TestRegistrationsEndpoint(t *testing.T) {
	suite.Run(t, new(RegistrationTestSuite))
}

func (suite *RegistrationTestSuite) TestEmptyBody() {
	body := []byte(`{}`)
	req := httptest.NewRequest("POST", "http://foobar/registrations", bytes.NewReader(body))
	RegistrationHandler(suite.rec, req)
	suite.Equal(http.StatusBadRequest, suite.rec.Result().StatusCode)
}

func (suite *RegistrationTestSuite) TestNoBody() {
	body := []byte(``)
	req := httptest.NewRequest("POST", "http://foobar/registrations", bytes.NewReader(body))
	RegistrationHandler(suite.rec, req)
	suite.Equal(http.StatusBadRequest, suite.rec.Result().StatusCode)
}

func (suite *RegistrationTestSuite) TestBadBody() {
	body := []byte(`{`)
	req := httptest.NewRequest("POST", "http://foobar/registrations", bytes.NewReader(body))
	RegistrationHandler(suite.rec, req)
	suite.Equal(http.StatusBadRequest, suite.rec.Result().StatusCode)
}

func (suite *RegistrationTestSuite) TestNoToken() {
	body := []byte(`{"uid": "abc1234"}`)
	req := httptest.NewRequest("POST", "http://foobar/registrations", bytes.NewReader(body))
	RegistrationHandler(suite.rec, req)
	suite.Equal(http.StatusBadRequest, suite.rec.Result().StatusCode)
}

func (suite *RegistrationTestSuite) TestNonBearerToken() {
	body := []byte(`{"uid": "abc1234"}`)
	req := httptest.NewRequest("POST", "http://foobar/registrations", bytes.NewReader(body))
	req.Header.Add("Authorization", "TrustmeBro 1234")
	RegistrationHandler(suite.rec, req)
	suite.Equal(http.StatusBadRequest, suite.rec.Result().StatusCode)
}

func (suite *RegistrationTestSuite) TestBadToken() {
	body := []byte(`{"uid": "abc1234"}`)
	req := httptest.NewRequest("POST", "http://foobar/registrations", bytes.NewReader(body))
	req.Header.Add("Authorization", "TrustmeBro1234")
	RegistrationHandler(suite.rec, req)
	suite.Equal(http.StatusBadRequest, suite.rec.Result().StatusCode)
}

func (suite *RegistrationTestSuite) TestExistingRegistration() {
	_, err := suite.store.Create(&store.Registration{UID: "abc1234"})
	suite.Nil(err)

	body := []byte(`{"uid": "abc1234"}`)
	req := httptest.NewRequest("POST", "http://foobar/registrations", bytes.NewReader(body))
	req.Header.Add("Authorization", "Bearer 1234")

	RegistrationHandler(suite.rec, req)

	suite.Equal(http.StatusConflict, suite.rec.Result().StatusCode)
}

func (suite *RegistrationTestSuite) TestSuccessfulRegistration() {
	body := []byte(`{"uid": "abc1234"}`)
	req := httptest.NewRequest("POST", "http://foobar/registrations", bytes.NewReader(body))
	req.Header.Add("Authorization", "Bearer 1234")

	RegistrationHandler(suite.rec, req)
	suite.Equal(http.StatusCreated, suite.rec.Result().StatusCode)
}
