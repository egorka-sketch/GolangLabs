package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
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

func (ur *UserRepo) GetAll(ctx context.Context, filter UserFilter) ([]User, error) {
	query := "SELECT id,name,birth_day FROM users"
	var conds []string
	var args []any
	if filter.Name != "" {
		args = append(args, "%"+filter.Name+"%")
		conds = append(conds, fmt.Sprintf("name ILIKE $%d", len(args)))
	}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if filter.MinAge >= 0 {
		args = append(args, today.AddDate(-filter.MinAge, 0, 0))
		conds = append(conds, fmt.Sprintf("birth_day <= $%d", len(args)))
	}
	if filter.MaxAge >= 0 {
		args = append(args, today.AddDate(-(filter.MaxAge+1), 0, 0))
		conds = append(conds, fmt.Sprintf("birth_day > $%d", len(args)))
	}
	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d", len(args))
	args = append(args, filter.Offset)
	query += fmt.Sprintf(" OFFSET $%d", len(args))
	rows, err := ur.pool.Query(ctx, query, args...)
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

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания : %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("нет соединения с базой: %w", err)
	}
	return pool, nil
}
