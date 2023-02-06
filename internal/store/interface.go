package store

type Store interface {
	All() ([]Registration, error)
	// find a single registration using orgID OR uid, either one will work.
	Find(orgID, uid string) (*Registration, error)
	Create(r *Registration) (string, error)
	Update(r *Registration, update *RegistrationUpdate) error
	Delete(orgID, uid string) error
}
