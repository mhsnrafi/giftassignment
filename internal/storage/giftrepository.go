package storage

import (
	"context"
	"database/sql"
)

type GiftStorage interface {
}

type GiftRepository struct {
	db *sql.DB
}

func NewGiftRepository(db *sql.DB) *GiftRepository {
	return &GiftRepository{db}
}

func (gr *GiftRepository) GetAllAvailableGifts(ctx context.Context) ([]Gift, error) {
	rows, err := gr.db.QueryContext(ctx, "SELECT name, categories, istaken FROM gifts WHERE istaken = false ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gifts []Gift
	for rows.Next() {
		var gift Gift
		if err := rows.Scan(&gift.Name, &gift.Categories, &gift.IsTaken); err != nil {
			return nil, err
		}
		gifts = append(gifts, gift)
	}
	return gifts, nil
}

func (gr *GiftRepository) UpdateGiftTaken(ctx context.Context, giftName string) error {
	_, err := gr.db.ExecContext(ctx, "UPDATE gifts SET istaken = true WHERE name = $1", giftName)
	return err
}

func (gr *GiftRepository) GetRandomAvailableGift(ctx context.Context) (Gift, error) {
	// Initialize a Gift variable to hold the results
	var gift Gift

	// Query to get a random available gift
	query := "SELECT * FROM gifts WHERE isTaken = false ORDER BY RANDOM() LIMIT 1;"

	// Execute the query
	row := gr.db.QueryRowContext(ctx, query)

	// Scan the results into the Gift variable
	err := row.Scan(&gift.ID, &gift.Name, &gift.Categories, &gift.IsTaken, &gift.Returned)
	if err != nil {
		return Gift{}, err
	}

	return gift, nil
}

func (repo *GiftRepository) UpdateGiftReturnStatus(ctx context.Context, giftName string, returned bool) error {
	// Create a prepared statement for executing an SQL update
	stmt, err := repo.db.PrepareContext(ctx, "UPDATE gifts SET returned = $1 WHERE name = $2")
	if err != nil {
		return err
	}

	// Execute the SQL update
	_, err = stmt.ExecContext(ctx, returned, giftName)
	if err != nil {
		return err
	}

	return nil
}
