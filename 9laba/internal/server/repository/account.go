package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Account struct {
	ID       int
	Login    string
	Password string
	Role     string
}

var ErrLoginTaken = errors.New("логин занят")

type AccountRepo interface {
	Create(ctx context.Context, login string, password string) (Account, error)
	Get(ctx context.Context, login string) (Account, error)
}

type AccRepo struct {
	pool *pgxpool.Pool
}

var _ AccountRepo = (*AccRepo)(nil)

func NewAccRepo(pool *pgxpool.Pool) *AccRepo {
	return &AccRepo{pool: pool}
}

func (r *AccRepo) Create(ctx context.Context, login string, password string) (Account, error) {
	var acc Account
	err := r.pool.QueryRow(ctx, "INSERT INTO accounts (login, password_hash) VALUES ($1, $2) RETURNING id, login, password_hash, role", login, password).Scan(&acc.ID, &acc.Login, &acc.Password, &acc.Role)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.Code == "23505" {
				return Account{}, ErrLoginTaken
			}
		}
		return Account{}, fmt.Errorf("создание аккаунта: %w", err)
	}
	return acc, nil
}

func (r *AccRepo) Get(ctx context.Context, login string) (Account, error) {
	var acc Account
	err := r.pool.QueryRow(ctx, "SELECT id, login, password_hash, role FROM accounts WHERE login = $1", login).Scan(&acc.ID, &acc.Login, &acc.Password, &acc.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Account{}, ErrNotFound
		}
		return Account{}, fmt.Errorf("поиск аккаунта: %w", err)
	}
	return acc, nil
}
