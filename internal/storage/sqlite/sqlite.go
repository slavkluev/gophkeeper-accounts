package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"accounts/internal/domain/models"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetAll(ctx context.Context, userUID uint64) ([]models.Account, error) {
	const op = "storage.sqlite.GetAll"

	stmt, err := s.db.Prepare("SELECT id, login, pass, info FROM accounts WHERE user_uid = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userUID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var accounts []models.Account

	for rows.Next() {
		var account models.Account
		err = rows.Scan(&account.ID, &account.Login, &account.Pass, &account.Info)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *Storage) SaveAccount(ctx context.Context, login string, pass string, info string, userUID uint64) (uint64, error) {
	const op = "storage.sqlite.SaveAccount"

	stmt, err := s.db.Prepare("INSERT INTO accounts(login, pass, info, user_uid) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, login, pass, info, userUID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return uint64(id), nil
}

func (s *Storage) UpdateAccount(
	ctx context.Context,
	id uint64,
	login string,
	pass string,
	info string,
	userUID uint64,
) error {
	const op = "storage.sqlite.UpdateAccount"

	stmt, err := s.db.Prepare("UPDATE accounts SET login = ?, pass = ?, info = ? WHERE id = ? AND user_uid = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, login, pass, info, id, userUID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
