package store

/*
Registration represents an instance of a satellite that is registered via:
- OrgID; comes from keycloak
- Uid; the CN from the satellite certificate

ID is a generated UUID
Extra is just a jsonb column if we want to store some extra metadata someday
*/
type Registration struct {
	ID    string
	OrgID string
	UID   string
	Extra map[string]interface{}
}

type RegistrationUpdate struct {
	Extra *map[string]interface{}
}
