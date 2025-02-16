package repository

import (
	"database/sql"
	"golang.org/x/net/context"
)

type TransactionRepository struct {
	db *sql.DB
}

func (t TransactionRepository) SendMoney(c context.Context, senderId int, receiver string, amount int32) error {
	tx, err := t.db.BeginTx(c, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	query := `
        WITH sender_update AS (
           UPDATE users
           SET coins = coins - $1
           WHERE id = $2 AND coins >= $1
           RETURNING id, coins
       )
       , receiver_update AS (
           UPDATE users
           SET coins = coins + $1
           WHERE username = $3 
           AND EXISTS (SELECT 1 FROM sender_update)
           RETURNING id
       )
       , transaction_insert AS (
           INSERT INTO transactions (sender_id, receiver_id, amount)
           SELECT sender_update.id, receiver_update.id, $1
           FROM sender_update, receiver_update
       )
       SELECT coins FROM sender_update;
    `

	var balance int32

	row := tx.QueryRowContext(c, query, amount, senderId, receiver)
	if err = row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	if err = row.Scan(&balance); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}
