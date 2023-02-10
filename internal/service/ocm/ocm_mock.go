package ocm

import (
	"context"
	"fmt"

	"github.com/redhatinsights/mbop/internal/models"
)

type OcmSDKMock struct{}

func (ocm *OcmSDKMock) InitSdkConnection(ctx context.Context) error {
	return nil
}

func (ocm *OcmSDKMock) GetUsers(usernames models.UserBody, q models.UserQuery) (models.Users, error) {
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

func (ocm *OcmSDKMock) IsOrgAdmin(id string) (bool, error) {
	if id == "23456" {
		return false, nil
	}

	if id == "errorTest" {
		return false, fmt.Errorf("error retrieving Role Bindings")
	}

	return true, nil
}

func (ocm *OcmSDKMock) CloseSdkConnection() {
	// nil
}
