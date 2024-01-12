package repository

import (
	"context"
	"database/sql"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository/transaction"
)

const coupleTableName = "couples"

type CoupleRepository interface {
	Has(usersDirection model.Direction) (bool, error)
	Add(ctx context.Context, couple *model.Couple) error
}

type coupleDBRepository struct {
	db *sql.DB
}

func NewCoupleDBRepository(db *sql.DB) CoupleRepository {
	return &coupleDBRepository{db}
}

func (cr *coupleDBRepository) Has(usersDirection model.Direction) (bool, error) {
	if err := usersDirection.Validate(); err != nil {
		return false, err
	}

	sb := newSelectBuilder()
	query, args := sb.
		Select("1").
		From(coupleTableName).
		Where(
			sb.Equal("from_user_id", usersDirection.FromID), sb.Equal("to_user_id", usersDirection.ToID),
			sb.Or(
				sb.Equal("from_user_id", usersDirection.ToID), sb.Equal("to_user_id", usersDirection.FromID),
			),
		).
		Limit(1).
		Build()

	res, err := cr.db.Exec(query, args...)
	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()

	return ra > 0, err
}

func (cr *coupleDBRepository) Add(ctx context.Context, couple *model.Couple) error {
	couple.BeforeAdd()

	if err := couple.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().
		InsertInto(coupleTableName).
		Cols("created_at", "updated_at", "from_user_id", "to_user_id").
		Values(couple.CreatedAt, couple.UpdatedAt, couple.UsersDirection.FromID, couple.UsersDirection.ToID).
		Build()

	conn := transaction.TxOrDB(ctx, cr.db)
	_, err := conn.Exec(query, args...)

	return err
}
