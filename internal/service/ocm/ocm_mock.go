package ocm

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/google/uuid"
	"github.com/redhatinsights/mbop/internal/models"
)

type SDKMock struct{}

func (ocm *SDKMock) InitSdkConnection(ctx context.Context) error {
	return nil
}

func (ocm *SDKMock) GetUsers(u models.UserBody, q models.UserQuery) (models.Users, error) {
	var users models.Users

	if u.Usernames == nil {
		return users, nil
	}

	if u.Usernames[0] == "errorTest" {
		return users, fmt.Errorf("internal AMS Error")
	}

	for _, user := range u.Usernames {
		users.AddUser(models.User{
			Username:      user,
			ID:            uuid.New().String(),
			Email:         "lub@dub.com",
			FirstName:     "test",
			LastName:      "case",
			AddressString: "https://usersTest.com",
			IsActive:      true,
			IsInternal:    true,
			Locale:        "en_US",
			OrgID:         strconv.Itoa(rand.Intn(999999 - 100000)),
			DisplayName:   "FedRAMP" + strconv.Itoa(rand.Intn(90-0)),
			Type:          "User",
		})
	}

	return users, nil
}

func (ocm *SDKMock) GetOrgAdmin(users []models.User) (models.OrgAdminResponse, error) {
	response := models.OrgAdminResponse{}

	if users[0].ID == "23456" {
		return response, nil
	}

	if users[0].ID == "errorTest" {
		return response, fmt.Errorf("error retrieving Role Bindings")
	}

	for _, user := range users {
		response[user.ID] = models.OrgAdmin{
			ID:         user.ID,
			IsOrgAdmin: true,
		}
	}

	return response, nil
}

func (ocm *SDKMock) CloseSdkConnection() {
	// nil
}
