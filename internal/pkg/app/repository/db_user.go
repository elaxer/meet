package repository

import (
	"context"
	"database/sql"
	"errors"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository/transaction"
)

const userTableName = "users"

type UserRepository interface {
	GetByLogin(login string) (*model.User, error)
	HasByLogin(login string) (bool, error)
	GetByTgID(id int64) (*model.User, error)
	HasByTgID(id int64) (bool, error)
	Add(ctx context.Context, user *model.User) error
	Update(user *model.User) error
	Remove(user *model.User) error
}

type userDBRepository struct {
	db *sql.DB
}

func NewUserDBRepository(db *sql.DB) UserRepository {
	return &userDBRepository{db}
}

func (ur *userDBRepository) GetByLogin(login string) (*model.User, error) {
	sb := newSelectBuilder()

	query, args := sb.Select("*").From(userTableName).Where(sb.Equal("login", login)).Limit(1).Build()

	u := new(model.User)

	row := ur.db.QueryRow(query, args...)
	err := row.Scan(u.GetFieldPointers()...)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	return u, err
}

func (ur *userDBRepository) HasByLogin(login string) (bool, error) {
	sb := newSelectBuilder()

	query, args := sb.Select("*").From(userTableName).Where(sb.Equal("login", login)).Limit(1).Build()

	res, err := ur.db.Exec(query, args...)
	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()

	return ra > 0, err
}

func (ur *userDBRepository) GetByTgID(id int64) (*model.User, error) {
	sb := newSelectBuilder()

	query, args := sb.Select("*").From(userTableName).Where(sb.Equal("tg_id", id)).Limit(1).Build()

	u := new(model.User)

	row := ur.db.QueryRow(query, args...)
	err := row.Scan(u.GetFieldPointers()...)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	return u, err
}

func (ur *userDBRepository) HasByTgID(id int64) (bool, error) {
	sb := newSelectBuilder()

	query, args := sb.Select("*").From(userTableName).Where(sb.Equal("tg_id", id)).Limit(1).Build()

	res, err := ur.db.Exec(query, args...)
	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()

	return ra > 0, err
}

func (ur *userDBRepository) Add(ctx context.Context, user *model.User) error {
	user.BeforeAdd()

	if err := user.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().
		InsertInto(userTableName).
		Cols("created_at", "login", "password_hash", "is_blocked", "tg_id").
		Values(user.CreatedAt, user.Login, user.PasswordHash, user.IsBlocked, user.TgID).
		SQL("RETURNING id").
		Build()

	var id int

	conn := transaction.TxOrDB(ctx, ur.db)
	row := conn.QueryRow(query, args...)
	if err := row.Scan(&id); err != nil {
		return err
	}

	user.ID = id

	return nil
}

func (ur *userDBRepository) Update(user *model.User) error {
	user.BeforeUpdate()

	if err := user.Validate(); err != nil {
		return err
	}

	ub := newUpdateBuilder()
	query, args := ub.Update(userTableName).
		Set(
			ub.Assign("updated_at", user.UpdatedAt),
			ub.Assign("password_hash", user.PasswordHash),
			ub.Assign("is_blocked", user.IsBlocked),
			ub.Assign("tg_id", user.TgID),
		).
		Where(ub.Equal("id", user.ID)).
		Build()

	_, err := ur.db.Exec(query, args...)

	return err
}

func (ur *userDBRepository) Remove(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	db := newDeleteBuilder()
	query, args := db.
		DeleteFrom(userTableName).
		Where(db.Equal("id", user.ID)).
		Build()

	_, err := ur.db.Exec(query, args...)

	return err
}
