package users

import (
	"fmt"

	sdk "github.com/openshift-online/ocm-sdk-go"
	v1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	"github.com/redhatinsights/mbop/internal/models"
)

func ResponseToUsers(response *v1.AccountsListResponse, connection *sdk.Connection) (models.Users, error) {
	users := models.Users{}
	items := response.Items().Slice()

	for i := range items {
		orgAdmin, err := isOrgAdmin(connection, items[i].ID())

		if err != nil {
			return users, err
		}

		users.AddUser(models.User{
			Username:      items[i].Username(),
			ID:            items[i].ID(),
			Email:         items[i].Email(),
			FirstName:     items[i].FirstName(),
			LastName:      items[i].LastName(),
			AddressString: items[i].HREF(),
			IsActive:      true,
			IsOrgAdmin:    orgAdmin,
			IsInternal:    true,
			Locale:        "en_US",
			OrgID:         items[i].Organization().ExternalID(),
			DisplayName:   items[i].Organization().Name(),
			Type:          items[i].Kind(),
		})
	}

	return users, nil
}

func CreateSearchString(usernames models.UserBody) string {
	search := ""

	for i := range usernames.Users {
		if i > 0 {
			search += " and "
		}

		search += fmt.Sprint("username='%s'", usernames.Users[i])
	}

	return search
}

func isOrgAdmin(connection *sdk.Connection, id string) (bool, error) {
	search := fmt.Sprintf("account.id='%s' and role.id='OrganizationAdmin'", id)

	collection := connection.AccountsMgmt().V1().RoleBindings()
	roleBindings, err := collection.List().Search(search).Send()

	if err != nil {
		return false, err
	}

	if roleBindings.Items().Empty() {
		return false, err
	}

	return true, err
}
