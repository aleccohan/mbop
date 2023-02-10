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
	IsOrgAdmin(id string) (bool, error)
}

func NewOcmClient() (OCM, error) {
	var client OCM

	mod := os.Getenv("OCM_MODULE")
	switch mod {
	case "aws":
		client = &OcmSDK{}
	case "mock":
		client = &OcmSDKMock{}
	default:
		return nil, fmt.Errorf("unsupported ocm module %q", mod)
	}

	return client, nil
}
