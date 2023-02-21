package models

type Users struct {
	Users []User `json:"users,omitempty"`
}

type UserV3Responses struct {
	Responses []UserV3Response `json:"responses,omitempty"`
}

type User struct {
	Username      string `json:"username"`
	ID            string `json:"id"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	AccountNumber string `json:"account_number"`
	AddressString string `json:"address_string"`
	IsActive      bool   `json:"is_active"`
	IsOrgAdmin    bool   `json:"is_org_admin"`
	IsInternal    bool   `json:"is_internal"`
	Locale        string `json:"locale"`
	OrgID         string `json:"org_id"`
	DisplayName   string `json:"display_name"`
	Entitlements  string `json:"entitlements"`
	Type          string `json:"type"`
}

type UserV3Response struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	IsActive   bool   `json:"is_active"`
	IsOrgAdmin bool   `json:"is_org_admin"`
	IsInternal bool   `json:"is_internal"`
	Locale     string `json:"locale"`
}

type UserV1Query struct {
	SortOrder string `json:"sortOrder"`
	QueryBy   string `json:"queryBy"`
}

type UserV3Query struct {
	SortOrder string `json:"sortOrder"`
	AdminOnly bool   `json:"admin_only"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

type UserBody struct {
	Users []string `json:"users"`
}

type UsersByBody struct {
	PrimaryEmail        string `json:"primary_email"`
	EmailStartsWith     string `json:"email_starts_with"`
	PrincipalStartsWith string `json:"principal_starts_with"`
}

func (u *Users) AddUser(user User) {
	u.Users = append(u.Users, user)
}

func (u *Users) RemoveUser(index int) {
	u.Users = append(u.Users[:index], u.Users[index+1:]...)
}

func (r *UserV3Responses) AddV3Response(response UserV3Response) {
	r.Responses = append(r.Responses, response)
}
