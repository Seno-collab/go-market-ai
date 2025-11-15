package db

// import (
// 	"context"
// 	"errors"

// 	"github.com/jackc/pgx/v4/pgxpool"
// 	"github.com/jackc/pgx/v5"
// )

// type UserRepo struct {
// 	pool *pgxpool.Pool
// 	q    *sqlc.Queries
// }

// func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
// 	return &UserRepo{pool: pool, q: sqlc.New(pool)}
// }

// func (r *UserRepo) GetByID(id int64) (user.User, error) {
// 	u, err := r.q.GetUserByID(context.Background(), id)
// 	if errors.Is(err, pgx.ErrNoRows) {
// 		return user.User{}, user.ErrNotFound
// 	}
// 	if err != nil {
// 		return user.User{}, err
// 	}
// 	return user.User{ID: u.ID, Email: u.Email, Name: u.Name}, nil
// }

// func (r *UserRepo) GetByEmail(email string) (user.User, error) {
// 	u, err := r.q.GetUserByEmail(context.Background(), email)
// 	if errors.Is(err, pgx.ErrNoRows) {
// 		return user.User{}, user.ErrNotFound
// 	}
// 	if err != nil {
// 		return user.User{}, err
// 	}
// 	return user.User{ID: u.ID, Email: u.Email, Name: u.Name}, nil
// }

// func (r *UserRepo) Create(in user.User) (int64, error) {
// 	u, err := r.q.CreateUser(context.Background(), sqlc.CreateUserParams{
// 		Email: in.Email,
// 		Name:  in.Name,
// 	})
// 	if err != nil {
// 		// TODO: check unique_violation via pgerrcode
// 		return 0, err
// 	}
// 	return u.ID, nil
// }
