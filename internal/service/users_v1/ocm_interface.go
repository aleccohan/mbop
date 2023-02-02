package users

import (
	"context"
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/openshift-online/ocm-sdk-go/logging"
)

type OcmSDK struct{}

func (ocm *OcmSDK) InitSdkConnection(ctx context.Context) (*sdk.Connection, error) {
	// Create a logger that has the debug level enabled:
	logger, err := logging.NewGoLoggerBuilder().
		Debug(true).
		Build()

	if err != nil {
		return nil, err
	}

	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).

		// SA Auth:
		Client(os.Getenv("AWS_CLIENT_ID"), os.Getenv("AWS_CLIENT_SECRET")).

		// Offline Token Auth:
		// Tokens(<token>).

		// Oauth Token URL:
		TokenURL(os.Getenv("OAUTH_TOKEN_URL")).

		// Route to hit for AMS:
		URL(os.Getenv("AMS_URL")).

		// SA Scopes:
		Scopes(os.Getenv("COGNITO_SCOPE")).
		BuildContext(ctx)

	return connection, err
}

func (ocm *OcmSDK) CloseSdkConnection(connection *sdk.Connection) {
	defer connection.Close()
}
