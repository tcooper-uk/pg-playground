package simulator

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

var firstNames = []string{"Alice", "Bob", "Carol", "David", "Eva", "Frank", "Grace", "Henry"}
var lastNames = []string{"Smith", "Jones", "Williams", "Brown", "Taylor", "Davies", "Evans", "Wilson"}

func simulateChurn(ctx context.Context, pool *pgxpool.Pool, ids *lookupIDs) error {
	if rand.Intn(10) < 7 {
		return createCustomer(ctx, pool, ids)
	}
	return deleteInactiveCustomer(ctx, pool)
}

func createCustomer(ctx context.Context, pool *pgxpool.Pool, ids *lookupIDs) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	cityID := ids.randomCityID()
	first := firstNames[rand.Intn(len(firstNames))]
	last := lastNames[rand.Intn(len(lastNames))]

	var addressID int
	err = tx.QueryRow(ctx, `
		INSERT INTO address (address, district, city_id, phone)
		VALUES ($1, 'N/A', $2, '000-000-0000')
		RETURNING address_id`,
		fmt.Sprintf("%d Sim Street", rand.Intn(9999)+1), cityID,
	).Scan(&addressID)
	if err != nil {
		return fmt.Errorf("insert address: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO customer (store_id, first_name, last_name, email, address_id, activebool)
		VALUES ($1, $2, $3, $4, $5, true)`,
		ids.randomStoreID(),
		first, last,
		fmt.Sprintf("%s.%s.%d@sim.example", first, last, rand.Intn(9999)),
		addressID,
	)
	if err != nil {
		return fmt.Errorf("insert customer: %w", err)
	}

	return tx.Commit(ctx)
}

func deleteInactiveCustomer(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var customerID int
	err = tx.QueryRow(ctx, `
		SELECT customer_id FROM customer
		WHERE activebool = false
		ORDER BY random()
		LIMIT 1`).Scan(&customerID)
	if err != nil {
		return nil // no inactive customers to delete
	}

	if _, err = tx.Exec(ctx, `DELETE FROM payment WHERE customer_id = $1`, customerID); err != nil {
		return fmt.Errorf("delete payments: %w", err)
	}
	if _, err = tx.Exec(ctx, `DELETE FROM rental WHERE customer_id = $1`, customerID); err != nil {
		return fmt.Errorf("delete rentals: %w", err)
	}
	if _, err = tx.Exec(ctx, `DELETE FROM customer WHERE customer_id = $1`, customerID); err != nil {
		return fmt.Errorf("delete customer: %w", err)
	}

	return tx.Commit(ctx)
}
