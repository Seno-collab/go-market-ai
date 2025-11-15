package user

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Register(u User) (int64, error) {
	if err := u.Validate(); err != nil {
		return 0, err
	}
	// unique email
	if _, err := s.repo.GetByEmail(u.Email); err == nil {
		return 0, ErrConflict
	}
	return s.repo.Create(u)
}

func (s *Service) Detail(id int64) (User, error) {
	return s.repo.GetByID(id)
}
