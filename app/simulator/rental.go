package simulator

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func simulateRental(ctx context.Context, pool *pgxpool.Pool, ids *lookupIDs) error {
	customerID := ids.randomCustomerID()
	inventoryID := ids.randomInventoryID()
	staffID := ids.randomStaffID()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var rentalID int
	err = tx.QueryRow(ctx, `
		INSERT INTO rental (rental_date, inventory_id, customer_id, staff_id)
		VALUES (now(), $1, $2, $3)
		RETURNING rental_id`,
		inventoryID, customerID, staffID,
	).Scan(&rentalID)
	if err != nil {
		return fmt.Errorf("insert rental: %w", err)
	}

	var amount float64
	err = tx.QueryRow(ctx, `
		SELECT f.rental_rate
		FROM inventory i
		JOIN film f ON f.film_id = i.film_id
		WHERE i.inventory_id = $1`, inventoryID,
	).Scan(&amount)
	if err != nil {
		return fmt.Errorf("fetch rental rate: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO payment (customer_id, staff_id, rental_id, amount, payment_date)
		VALUES ($1, $2, $3, $4, now())`,
		customerID, staffID, rentalID, amount,
	)
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}

	return tx.Commit(ctx)
}
