package ocm

import (
	"context"
	"fmt"

	sdk "github.com/openshift-online/ocm-sdk-go"
	v1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	"github.com/openshift-online/ocm-sdk-go/logging"
	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/models"
)

type SDK struct {
	client *sdk.Connection
}

func (ocm *SDK) InitSdkConnection(ctx context.Context) error {
	// Create a logger that has the debug level enabled:
	logger, err := logging.NewGoLoggerBuilder().
		Debug(true).
		Build()

	if err != nil {
		return err
	}

	ocm.client, err = sdk.NewConnectionBuilder().
		Logger(logger).

		// SA Auth:
		Client(config.Get().CognitoAppClientID, config.Get().CognitoAppClientSecret).

		// Offline Token Auth:
		// Tokens(<token>).

		// Oauth Token URL:
		TokenURL(config.Get().OauthTokenURL).

		// Route to hit for AMS:
		URL(config.Get().AmsURL).

		// SA Scopes:
		Scopes(config.Get().CognitoScope).
		BuildContext(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (ocm *SDK) GetUsers(usernames models.UserBody, q models.UserQuery) (models.Users, error) {
	search := createSearchString(usernames)
	collection := ocm.client.AccountsMgmt().V1().Accounts().List().Search(search)

	collection = collection.Order(createQueryOrder(q))

	users := models.Users{}
	usersResponse, err := collection.Send()
	if err != nil {
		return users, err
	}

	if usersResponse.Items().Empty() {
		return users, err
	}

	users = responseToUsers(usersResponse)

	return users, err
}

func (ocm *SDK) GetOrgAdmin(u []models.User) (models.OrgAdminResponse, error) {
	search := createOrgAdminSearchString(u)

	collection := ocm.client.AccountsMgmt().V1().RoleBindings()
	roleBindings, err := collection.List().Search(search).Send()

	orgAdminResponse := models.OrgAdminResponse{}
	if err != nil {
		return orgAdminResponse, err
	}

	if roleBindings.Items().Empty() {
		return orgAdminResponse, err
	}

	bindingSlice := roleBindings.Items().Slice()
	for _, binding := range bindingSlice {
		orgAdminResponse[binding.Account().ID()] = models.OrgAdmin{
			ID:         binding.Account().ID(),
			IsOrgAdmin: true,
		}
	}

	return orgAdminResponse, err
}

func (ocm *SDK) CloseSdkConnection() {
	ocm.client.Close()
}

func responseToUsers(response *v1.AccountsListResponse) models.Users {
	users := models.Users{}
	items := response.Items().Slice()

	for i := range items {
		users.AddUser(models.User{
			Username:      items[i].Username(),
			ID:            items[i].ID(),
			Email:         items[i].Email(),
			FirstName:     items[i].FirstName(),
			LastName:      items[i].LastName(),
			AddressString: items[i].HREF(),
			IsActive:      true,
			IsInternal:    true,
			Locale:        "en_US",
			OrgID:         items[i].Organization().ID(),
			DisplayName:   items[i].Organization().Name(),
			Type:          items[i].Kind(),
		})
	}

	return users
}

func createSearchString(usernames models.UserBody) string {
	search := ""

	for i := range usernames.Users {
		if i > 0 {
			search += " or "
		}

		search += fmt.Sprintf("username='%s'", usernames.Users[i])
	}

	return search
}

func createOrgAdminSearchString(users []models.User) string {
	search := ""

	for i := range users {
		if i > 0 {
			search += " or "
		}

		search += fmt.Sprintf("account.id='%s' and role.id='OrganizationAdmin'", users[i].ID)
	}

	return search
}

func createQueryOrder(q models.UserQuery) string {
	order := ""

	if q.QueryBy != "" {
		order += q.QueryBy
	}

	if q.SortOrder != "" {
		order += fmt.Sprint(" " + q.SortOrder)
	}

	return order
}
