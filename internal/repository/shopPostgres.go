package repository

import (
	"database/sql"
	"golang.org/x/net/context"
)

type ShopRepository struct {
	db *sql.DB
}

func (s ShopRepository) BuyItem(c context.Context, userId int, item string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `WITH merch_data AS (
    	SELECT id AS merch_id, price FROM merch WHERE name = $1
	),
     	updated_user AS (
			UPDATE users
			SET coins = coins - (SELECT price FROM merch_data)
			WHERE id = $2 AND coins >= (SELECT price FROM merch_data)
			RETURNING id
	)
	INSERT INTO purchases (user_id, merch_id)
	SELECT updated_user.id, merch_data.merch_id FROM updated_user, merch_data
	RETURNING user_id;`
	row := tx.QueryRowContext(c, query, item, userId)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return err
	}

	var tempId int
	if err := row.Scan(&tempId); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
func NewShopRepository(db *sql.DB) *ShopRepository {
	return &ShopRepository{db: db}

}
