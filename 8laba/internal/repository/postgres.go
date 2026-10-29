package repository

import "github.com/jackc/pgx/v5/pgxpool"

type UserRepo struct {
	pool *pgxpool.Pool
}

//var _ UserRepository = (*UserRepo)(nil)

func NewRepo(pool *pgxpool.Pool) *UserRepo {
	db := &UserRepo{
		pool: pool,
	}
	return db
}
