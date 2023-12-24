package repository

import (
	"database/sql"
	"meet/internal/pkg/app"
	"meet/internal/pkg/app/model"

	"github.com/huandu/go-sqlbuilder"
)

const userTableName = "users"

type UserRepository interface {
	GetByLogin(login string) (*model.User, error)
	HasByLogin(login string) (bool, error)
	Add(user *model.User) error
	Update(user *model.User) error
	Remove(user *model.User) error
}

type userDBRepository struct {
	dbRepository
}

func newUserRepository(db *sql.DB) UserRepository {
	ur := new(userDBRepository)
	ur.dbRepository.db = db

	return ur
}

func (ur *userDBRepository) GetByLogin(login string) (*model.User, error) {
	sb := ur.createSelectBuilder()

	sql, args := sb.Select("*").From(userTableName).Where(sb.Equal("login", login)).Limit(1).Build()

	u := new(model.User)

	row := ur.db.QueryRow(sql, args...)
	err := row.Scan(u.GetFieldPointers()...)

	return u, err
}

func (ur *userDBRepository) HasByLogin(login string) (bool, error) {
	sb := ur.createSelectBuilder()

	sql, args := sb.Select("*").From(userTableName).Where(sb.Equal("login", login)).Limit(1).Build()

	res, err := ur.db.Exec(sql, args...)
	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()

	return ra > 0, err
}

func (ur *userDBRepository) Add(user *model.User) error {
	user.BeforeAdd()

	if err := user.Validate(); err != nil {
		return err
	}

	ib := sqlbuilder.NewInsertBuilder()
	sql, args := ib.InsertInto(userTableName).
		Cols("created_at", "login", "password_hash", "is_blocked").
		Values(user.CreatedAt, user.Login, user.PasswordHash, user.IsBlocked).
		BuildWithFlavor(app.SQLBuilderFlavor)

	_, err := ur.db.Exec(sql, args...)

	return err
}

func (ur *userDBRepository) Update(user *model.User) error {
	user.BeforeUpdate()

	if err := user.Validate(); err != nil {
		return err
	}

	ub := sqlbuilder.NewUpdateBuilder()
	sql, args := ub.Update(userTableName).
		Set(
			ub.Assign("updated_at", user.UpdatedAt),
			ub.Assign("password_hash", user.PasswordHash),
			ub.Assign("is_blocked", user.IsBlocked),
		).
		Where(ub.Equal("id", user.ID)).
		BuildWithFlavor(app.SQLBuilderFlavor)

	_, err := ur.db.Exec(sql, args...)

	return err
}

func (ur *userDBRepository) Remove(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	db := sqlbuilder.NewDeleteBuilder()
	query, args := db.
		DeleteFrom(userTableName).
		Where(db.Equal("id", user.ID)).
		BuildWithFlavor(app.SQLBuilderFlavor)

	_, err := ur.db.Exec(query, args...)

	return err
}
