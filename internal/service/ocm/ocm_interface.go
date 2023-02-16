package ocm

import (
	"context"
	"fmt"

	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/models"
)

type OCM interface {
	InitSdkConnection(ctx context.Context) error
	CloseSdkConnection()
	GetUsers(users models.UserBody, q models.UserQuery) (models.Users, error)
	GetOrgAdmin([]models.User) (models.OrgAdminResponse, error)
}

// re-declaring ams constant here to avoid circular module importing
const amsModule = "ams"

func NewOcmClient() (OCM, error) {
	var client OCM

	switch config.Get().UsersModule {
	case amsModule:
		client = &SDK{}
	case "mock":
		client = &SDKMock{}
	default:
		return nil, fmt.Errorf("unsupported users module %q", config.Get().UsersModule)
	}

	return client, nil
}
