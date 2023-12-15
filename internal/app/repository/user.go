package repository

import (
	"database/sql"
	"meet/internal/app"
	"meet/internal/app/model"

	"github.com/huandu/go-sqlbuilder"
)

const userTableName = "users"

type UserRepository interface {
	GetByLogin(login string) (*model.User, error)
	HasByLogin(login string) (bool, error)
	Add(user *model.User) error
	Update(user *model.User) error
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

	u := model.NewUser()

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
	if err := user.Validate(); err != nil {
		return err
	}
	hbl, err := ur.HasByLogin(user.Login)
	if err != nil {
		return err
	}
	if hbl {
		return ErrDuplicate
	}

	ib := sqlbuilder.NewInsertBuilder()
	sql, args := ib.InsertInto(userTableName).
		Cols("login", "password_hash", "is_blocked").
		Values(user.Login, user.PasswordHash, user.IsBlocked).
		BuildWithFlavor(app.SQLBuilderFlavor)

	_, err = ur.db.Exec(sql, args...)

	return err
}

func (ur *userDBRepository) Update(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	user.BeforeUpdate()

	ub := sqlbuilder.NewUpdateBuilder()
	sql, args := ub.Update(userTableName).
		Set(
			ub.Assign("login", user.Login),
			ub.Assign("password_hash", user.PasswordHash),
			ub.Assign("is_blocked", user.IsBlocked),
		).
		Where(ub.Equal("id", user.ID)).
		BuildWithFlavor(app.SQLBuilderFlavor)

	_, err := ur.db.Exec(sql, args...)

	return err
}
