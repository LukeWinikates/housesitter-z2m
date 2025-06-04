package database

type Store interface {
	Find(id string) (*Schedule, error)
	All() []*Schedule
	SaveChanges(id string, s *Schedule) error
}
