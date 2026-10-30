package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

var _ UserRepository = (*UserRepo)(nil)

func NewRepo(pool *pgxpool.Pool) *UserRepo {
	db := &UserRepo{
		pool: pool,
	}
	return db
}

func (ur *UserRepo) GetAll(ctx context.Context) ([]User, error) {
	rows, err := ur.pool.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("ошибка вывода: %w", err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err = rows.Scan(&u.ID, &u.Name, &u.Birthday); err != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка вывода: %w", err)
	}
	return users, nil
}

func (ur *UserRepo) GetByID(ctx context.Context, id int) (User, error) {
	var u User

	err := ur.pool.QueryRow(ctx, "SELECT * FROM users WHERE id=$1", id).Scan(&u.ID, &u.Name, &u.Birthday)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("ошибка получения пользователя:%d %w", id, err)
	}
	return u, nil
}

func (ur *UserRepo) Create(ctx context.Context, user User) (User, error) {
	var u User
	err := ur.pool.QueryRow(ctx, "insert into users(name, birth_day) values ($1,$2) RETURNING id, name, birth_day", user.Name, user.Birthday).Scan(&u.ID, &u.Name, &u.Birthday)
	if err != nil {
		return User{}, fmt.Errorf("ошибка создания профиля %w", err)
	}
	return u, nil
}

func (ur *UserRepo) Update(ctx context.Context, user User) (User, error) {
	var u User
	err := ur.pool.QueryRow(ctx, "update users set name =  $1, birth_day = $2  where id=$3 RETURNING id, name, birth_day", user.Name, user.Birthday, user.ID).Scan(&u.ID, &u.Name, &u.Birthday)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("ошибка обновления данных пользователя %d %w", user.ID, err)
	}
	return u, nil
}

func (ur *UserRepo) Delete(ctx context.Context, id int) error {
	tag, err := ur.pool.Exec(ctx, "delete from users where id=$1", id)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя : %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
