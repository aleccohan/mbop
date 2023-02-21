package ocm

import (
	"context"
	"fmt"

	sdk "github.com/openshift-online/ocm-sdk-go"
	v1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	"github.com/openshift-online/ocm-sdk-go/logging"
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
		// Client(os.Getenv("COGNITO_APP_CLIENT_ID"), os.Getenv("COGNITO_APP_CLIENT_SECRET")).
		Client("6k2fo38r9l306k9l2t1ji428jo", "1f8ontqbalms06td5375trbc1g2rmgmficeo146u2s6odinr9b2q").

		// Offline Token Auth:
		// Tokens(<token>).

		// Oauth Token URL:
		// TokenURL(os.Getenv("OAUTH_TOKEN_URL")).
		TokenURL("https://ocm-ra-stage-domain.auth-fips.us-gov-west-1.amazoncognito.com/oauth2/token").

		// Route to hit for AMS:
		// URL(os.Getenv("AMS_URL")).
		URL("https://ocm-stage.rosa-nlb.appsrefrs01ugw1.p1.openshiftusgov.com").

		// SA Scopes:
		// Scopes(os.Getenv("COGNITO_SCOPE")).
		Scopes("ocm/InsightsServiceAccount").
		BuildContext(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (ocm *SDK) GetUsers(usernames models.UserBody, q models.UserV1Query) (models.Users, error) {
	search := createSearchString(usernames)
	collection := ocm.client.AccountsMgmt().V1().Accounts().List().Search(search)

	collection = collection.Order(createQueryOrder(q))

	users := models.Users{Users: []models.User{}}
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

func (ocm *SDK) GetAccountV3Users(orgID string, q models.UserV3Query) (models.Users, error) {
	search := createAccountsV3UsersSearchString(orgID)

	collection := ocm.client.AccountsMgmt().V1().Accounts().List().Search(search)

	collection = collection.Order(createV3QueryOrder(q))
	collection = collection.Size(q.Limit)
	collection = collection.Page(q.Offset)

	users := models.Users{Users: []models.User{}}
	AccountV3UsersResponse, err := collection.Send()
	if err != nil {
		return users, err
	}

	users = responseToUsers(AccountV3UsersResponse)

	return users, err
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

func usersToV3Response(u models.Users) models.UserV3Response {
	response := models.UserV3Response{}

	for _, user := range u.Users {
		response.ID = user.OrgID
		response.Email = user.Email
		response.Username = user.Username
		response.FirstName = user.FirstName
		response.LastName = user.LastName
		response.IsActive = user.IsActive
		response.IsInternal = user.IsInternal
		response.Locale = user.Locale
	}

	return response
}

func createSearchString(u models.UserBody) string {
	search := ""

	for i := range u.Users {
		if i > 0 {
			search += " or "
		}

		search += fmt.Sprintf("username='%s'", u.Users[i])
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

func createAccountsV3UsersSearchString(orgID string) string {
	return fmt.Sprintf("organization.id='%s'", orgID)
}

func createQueryOrder(q models.UserV1Query) string {
	order := ""

	if q.QueryBy != "" {
		order += q.QueryBy
	}

	if q.SortOrder != "" {
		order += fmt.Sprint(" " + q.SortOrder)
	}

	return order
}

func createV3QueryOrder(q models.UserV3Query) string {
	order := "organization.id"

	if q.SortOrder != "" {
		order += fmt.Sprint(" " + q.SortOrder)
	}

	return order
}
