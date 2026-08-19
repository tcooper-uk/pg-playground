package simulator

import (
	"context"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

func simulateRead(ctx context.Context, pool *pgxpool.Pool, ids *lookupIDs) error {
	switch rand.Intn(3) {
	case 0:
		return readFilm(ctx, pool, ids)
	case 1:
		return readActorFilmography(ctx, pool, ids)
	default:
		return readCustomerHistory(ctx, pool, ids)
	}
}

func readFilm(ctx context.Context, pool *pgxpool.Pool, ids *lookupIDs) error {
	rows, err := pool.Query(ctx, `
		SELECT f.title, f.rating, f.rental_rate, l.name
		FROM film f
		JOIN language l ON l.language_id = f.language_id
		WHERE f.film_id = $1`, ids.randomFilmID())
	if err != nil {
		return err
	}
	rows.Close()
	return rows.Err()
}

func readActorFilmography(ctx context.Context, pool *pgxpool.Pool, ids *lookupIDs) error {
	rows, err := pool.Query(ctx, `
		SELECT a.first_name, a.last_name, f.title
		FROM actor a
		JOIN film_actor fa ON fa.actor_id = a.actor_id
		JOIN film f ON f.film_id = fa.film_id
		WHERE a.actor_id = $1`, ids.randomActorID())
	if err != nil {
		return err
	}
	rows.Close()
	return rows.Err()
}

func readCustomerHistory(ctx context.Context, pool *pgxpool.Pool, ids *lookupIDs) error {
	rows, err := pool.Query(ctx, `
		SELECT r.rental_date, r.return_date, f.title, p.amount
		FROM rental r
		JOIN inventory i ON i.inventory_id = r.inventory_id
		JOIN film f ON f.film_id = i.film_id
		LEFT JOIN payment p ON p.rental_id = r.rental_id
		WHERE r.customer_id = $1
		ORDER BY r.rental_date DESC
		LIMIT 20`, ids.randomCustomerID())
	if err != nil {
		return err
	}
	rows.Close()
	return rows.Err()
}
