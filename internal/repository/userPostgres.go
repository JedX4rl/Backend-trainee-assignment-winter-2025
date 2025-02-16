package repository

import (
	"Backend-trainee-assignment-winter-2025/internal/models"
	"database/sql"
	"golang.org/x/net/context"
)

type UserRepository struct {
	db *sql.DB
}

func (u UserRepository) SignUp(c context.Context, username, password string) (*models.User, error) {
	tx, err := u.db.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, coins`
	row := tx.QueryRowContext(c, query, username, password)

	var user models.User
	if err := row.Scan(&user.Id, &user.Coins); err != nil {
		return nil, err
	}
	user.Password = password
	user.Username = username

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) GetUserByUsername(c context.Context, username string) (*models.User, error) {

	query := `SELECT * FROM users WHERE username = $1`
	row := u.db.QueryRowContext(c, query, username)

	if err := row.Err(); err != nil {
		return nil, err
	}

	var user models.User

	if err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Coins); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserRepository) GetInfo(c context.Context, userId int) (*models.InfoResponse, error) {
	query := `WITH u_transactions AS (
    SELECT t.sender_id AS user_id, u_sender.username AS other_user, t.amount, FALSE AS is_sent
    FROM transactions t
             JOIN users u_sender ON u_sender.id = t.sender_id
    WHERE t.receiver_id = $1
    UNION ALL
    SELECT t.receiver_id AS user_id, u_receiver.username AS other_user, t.amount, TRUE AS is_sent
    FROM transactions t
             JOIN users u_receiver ON u_receiver.id = t.receiver_id
    WHERE t.sender_id = $1
),
     u_inventory AS (
         SELECT p.user_id, m.name AS item, COUNT(*) AS quantity
         FROM purchases p
                  JOIN merch m ON p.merch_id = m.id
         WHERE p.user_id = $1
         GROUP BY p.user_id, m.name
     )
SELECT
    u.coins,
    t.amount, t.other_user, t.is_sent,
    i.item, i.quantity
FROM users u
         LEFT JOIN u_transactions t ON u.id = t.user_id
         LEFT JOIN users sender ON sender.id = t.user_id
         LEFT JOIN u_inventory i ON u.id = i.user_id
WHERE u.id = $1;`
	tx, err := u.db.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(c, query, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	defer rows.Close()

	var infoResponse models.InfoResponse

	var received, sent []models.TransactionInfo
	var inventory []models.InventoryItem

	for rows.Next() {
		var t models.TransactionInfo
		var i models.InventoryItem
		var otherUsername sql.NullString
		var amount sql.NullInt32
		var isSent sql.NullBool
		var itemType sql.NullString
		var quantity sql.NullInt64

		if err = rows.Scan(&infoResponse.Coins, &amount, &otherUsername, &isSent, &itemType, &quantity); err != nil {
			tx.Rollback()
			return nil, err
		}

		if amount.Valid && otherUsername.Valid && isSent.Valid {
			t.Username = otherUsername.String
			t.Amount = amount.Int32
			t.IsSent = isSent.Bool
			if isSent.Bool {
				sent = append(sent, t)
			} else {
				received = append(received, t)
			}
		}
		if itemType.Valid && quantity.Valid {
			i.Type = itemType.String
			i.Quantity = quantity.Int64
			inventory = append(inventory, i)
		}
	}
	infoResponse.Inventory = inventory
	infoResponse.CoinHistory.Received = received
	infoResponse.CoinHistory.Sent = sent
	return &infoResponse, tx.Commit()
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
