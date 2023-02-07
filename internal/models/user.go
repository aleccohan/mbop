package models

type Users struct {
	Users []User `json:"users,omitempty"`
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

type UserQuery struct {
	SortOrder string `json:"sortOrder"`
	QueryBy   string `json:"queryBy"`
}

func (u *Users) AddUser(user User) {
	u.Users = append(u.Users, user)
}
