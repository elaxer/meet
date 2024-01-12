package repository

import (
	"database/sql"
	"errors"
	"meet/internal/pkg/app/model"
)

const messageTableName = "messages"

type MessageRepository interface {
	Get(id int) (*model.Message, error)
	GetList(usersDirection model.Direction, limit, offset int) ([]*model.Message, error)
	UnreadCount(userID int) (int, error)
	Add(message *model.Message) error
	Update(message *model.Message) error
}

type messageDBRepository struct {
	db *sql.DB
}

func NewMessageDBRepository(db *sql.DB) MessageRepository {
	return &messageDBRepository{db}
}

func (mr *messageDBRepository) Get(id int) (*model.Message, error) {
	sb := newSelectBuilder()
	query, args := sb.Select("*").From(messageTableName).Where(sb.Equal("id", id)).Limit(1).Build()

	m := new(model.Message)
	row := mr.db.QueryRow(query, args...)
	if err := row.Scan(m.GetFieldPointers()...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return m, nil
}

func (mr *messageDBRepository) GetList(usersDirection model.Direction, limit, offset int) ([]*model.Message, error) {
	messages := make([]*model.Message, 0, limit)

	sb := newSelectBuilder()

	query, args := sb.
		Select("*").
		From(messageTableName).
		Where(sb.Or(
			sb.And(sb.Equal("from_user_id", usersDirection.FromID), sb.Equal("to_user_id", usersDirection.ToID)),
			sb.And(sb.Equal("from_user_id", usersDirection.ToID), sb.Equal("to_user_id", usersDirection.FromID)),
		)).
		OrderBy("created_at").
		Desc().
		Limit(limit).
		Offset(offset).
		Build()

	rows, err := mr.db.Query(query, args...)
	if err != nil {
		return messages[0:0], err
	}

	for rows.Next() {
		m := new(model.Message)

		if err := rows.Scan(m.GetFieldPointers()...); err != nil {
			return messages[0:0], err
		}

		messages = append(messages, m)
	}

	return messages, nil
}

func (mr *messageDBRepository) UnreadCount(userID int) (int, error) {
	sb := newSelectBuilder()

	query, args := sb.
		Select("*").
		From(messageTableName).
		Where(sb.Equal("to_user_id", userID)).
		Where(sb.Equal("is_readed", false)).
		Build()

	r, err := mr.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	ra, err := r.RowsAffected()

	return int(ra), err
}

func (mr *messageDBRepository) Add(message *model.Message) error {
	message.BeforeAdd()

	if err := message.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().
		InsertInto(messageTableName).
		Cols(
			"created_at",
			"from_user_id",
			"to_user_id",
			"text",
			"is_readed",
		).
		Values(
			message.CreatedAt,
			message.UsersDirection.FromID,
			message.UsersDirection.ToID,
			message.Text,
			message.IsReaded,
		).
		SQL("RETURNING id").
		Build()

	var id int
	row := mr.db.QueryRow(query, args...)
	if err := row.Scan(&id); err != nil {
		message.ID = id
	}

	return nil
}

func (mr *messageDBRepository) Update(message *model.Message) error {
	message.BeforeUpdate()

	if err := message.Validate(); err != nil {
		return err
	}

	ub := newUpdateBuilder()
	query, args := ub.
		Update(messageTableName).
		Set(
			ub.Assign("updated_at", message.UpdatedAt),
			ub.Assign("text", message.Text),
			ub.Assign("is_readed", message.IsReaded),
		).
		Build()

	_, err := mr.db.Exec(query, args...)

	return err
}
