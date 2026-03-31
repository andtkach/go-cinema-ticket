package movies

import "github.com/google/uuid"

type Service struct {
	store MovieStore
}

func NewService(store MovieStore) *Service {
	return &Service{store}
}

func (s *Service) List() ([]Movie, error) {
	return s.store.List()
}

func (s *Service) GetByID(id string) (Movie, error) {
	return s.store.GetByID(id)
}

func (s *Service) Create(m Movie) (Movie, error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return s.store.Create(m)
}

func (s *Service) Update(m Movie) (Movie, error) {
	return s.store.Update(m)
}

func (s *Service) Delete(id string) error {
	return s.store.Delete(id)
}
