package ocm

import (
	"context"
	"fmt"

	"github.com/redhatinsights/mbop/internal/models"
)

type SDKMock struct{}

func (ocm *SDKMock) InitSdkConnection(ctx context.Context) error {
	return nil
}

func (ocm *SDKMock) GetUsers(usernames models.UserBody, q models.UserQuery) (models.Users, error) {
	var users models.Users

	if usernames.Users == nil {
		return users, nil
	}

	if usernames.Users[0] == "errorTest" {
		return users, fmt.Errorf("internal AMS Error")
	}

	users = models.Users{
		Users: []models.User{
			{
				Username:      "test1",
				ID:            "12345",
				Email:         "lub@dub.com",
				FirstName:     "test",
				LastName:      "case",
				AddressString: "https://something.com",
				IsActive:      true,
				IsInternal:    true,
				Locale:        "en_US",
				OrgID:         "67890",
				DisplayName:   "FedRAMP1",
				Type:          "User",
			},
			{
				Username:      "test2",
				ID:            "23456",
				Email:         "lub@dub.com",
				FirstName:     "john",
				LastName:      "doe",
				AddressString: "https://something.com",
				IsActive:      true,
				IsInternal:    true,
				Locale:        "en_US",
				OrgID:         "78901",
				DisplayName:   "FedRAMP2",
				Type:          "User",
			},
		},
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

	response = models.OrgAdminResponse{
		"12345": models.OrgAdmin{
			ID:         "12345",
			IsOrgAdmin: true,
		},
	}

	return response, nil
}

func (ocm *SDKMock) CloseSdkConnection() {
	// nil
}
