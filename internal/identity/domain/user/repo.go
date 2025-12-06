package user

type Repository interface {
	GetByID(id int64) (User, error)
	GetByEmail(email string) (User, error)
	Create(u User) (int64, error)
}
