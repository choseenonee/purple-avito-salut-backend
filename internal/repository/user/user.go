package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"template/internal/model/entities"
	"template/internal/repository"
	"template/pkg/customerr"
)

type User struct {
	db *sqlx.DB
}

func InitUserRepo(db *sqlx.DB) repository.User {
	return User{
		db: db,
	}
}

func (u User) Create(ctx context.Context, userCreate entities.UserCreate) (int, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	var userID int

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userCreate.Password), 14)

	row := tx.QueryRowContext(ctx, `INSERT INTO users (name, hashed_password) VALUES ($1, $2) RETURNING id;`,
		userCreate.Name, hashedPassword)

	err = row.Scan(&userID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return userID, nil
}

func (u User) Get(ctx context.Context, userID int) (entities.User, error) {
	var user entities.User

	err := u.db.QueryRowContext(ctx, `SELECT id, name FROM users WHERE users.id = $1`,
		userID).Scan(&user.ID, &user.Name)

	if err != nil {
		return entities.User{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return user, nil
}

func (u User) GetHashedPassword(ctx context.Context, name string) (int, string, error) {
	var hashedPassword string
	var userID int

	err := u.db.QueryRowContext(ctx, `SELECT id, hashed_password FROM users WHERE users.name = $1`,
		name).Scan(&userID, &hashedPassword)

	if err != nil {
		return 0, "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return userID, hashedPassword, nil
}

func (u User) Delete(ctx context.Context, userID int) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM users WHERE users.id = $1;`, userID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}
	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}
	if count != 1 {
		err = errors.New("count error")
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}
