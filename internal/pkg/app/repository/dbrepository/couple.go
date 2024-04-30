package dbrepository

import (
	"context"
	"meet/internal/pkg/app/database"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
)

const coupleTableName = "couples"

type coupleRepository struct {
	conn database.Connection
}

func NewCoupleRepository(conn database.Connection) repository.CoupleRepository {
	return &coupleRepository{conn}
}

func (cr *coupleRepository) Has(usersDirection model.Direction) (bool, error) {
	if err := usersDirection.Validate(); err != nil {
		return false, err
	}

	sb := newSelectBuilder()
	query, args := sb.
		Select("1").
		From(coupleTableName).
		Where(
			sb.Or(
				sb.And(sb.Equal("from_user_id", usersDirection.FromID), sb.Equal("to_user_id", usersDirection.ToID)),
				sb.And(sb.Equal("from_user_id", usersDirection.ToID), sb.Equal("to_user_id", usersDirection.FromID)),
			),
		).
		Limit(1).
		Build()

	res, err := cr.conn.Exec(query, args...)
	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()

	return ra > 0, err
}

func (cr *coupleRepository) Add(ctx context.Context, couple *model.Couple) error {
	couple.BeforeAdd()

	if err := couple.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().
		InsertInto(coupleTableName).
		Cols("created_at", "updated_at", "from_user_id", "to_user_id").
		Values(couple.CreatedAt, couple.UpdatedAt, couple.UsersDirection.FromID, couple.UsersDirection.ToID).
		Build()

	conn := database.TxOrDB(ctx, cr.conn)
	_, err := conn.Exec(query, args...)

	return err
}
