package ocm

import (
	"context"
	"fmt"
	"os"

	"github.com/redhatinsights/mbop/internal/models"
)

type OCM interface {
	InitSdkConnection(ctx context.Context) error
	CloseSdkConnection()
	GetUsers(users models.UserBody, q models.UserQuery) (models.Users, error)
	GetOrgAdmin([]models.User) (models.OrgAdminResponse, error)
}

func NewOcmClient() (OCM, error) {
	var client OCM

	mod := os.Getenv("OCM_MODULE")
	switch mod {
	case "aws":
		client = &SDK{}
	case "mock":
		client = &SDKMock{}
	default:
		return nil, fmt.Errorf("unsupported ocm module %q", mod)
	}

	return client, nil
}
