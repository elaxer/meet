package fixture

import (
	"context"
	"database/sql"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/repository/transaction"
)

func LoadFixtures(
	db *sql.DB,
	userRepository repository.UserRepository,
	questionnaireRepository repository.QuestionnaireRepository,
	coupleRepository repository.CoupleRepository,
) error {
	ctx, tx, err := transaction.BeginTx(context.Background(), db)
	if err != nil {
		return err
	}

	for _, u := range Users {
		if err := userRepository.Add(ctx, u); err != nil {
			tx.Rollback()

			return err
		}
	}

	for u, q := range Questionnaires {
		q.UserID = u.ID

		if err := questionnaireRepository.Add(ctx, q); err != nil {
			tx.Rollback()

			return err
		}
	}

	for user, couples := range Couples {
		for _, couple := range couples {
			c := new(model.Couple)
			c.UsersDirection = model.Direction{FromID: user.ID, ToID: couple.ID}

			if err := coupleRepository.Add(ctx, c); err != nil {
				tx.Rollback()

				return err
			}
		}
	}

	err = tx.Commit()

	return err
}
