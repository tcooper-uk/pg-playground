package simulator

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func simulateReturn(ctx context.Context, pool *pgxpool.Pool) error {
	tag, err := pool.Exec(ctx, `
		UPDATE rental
		SET return_date = now()
		WHERE rental_id = (
			SELECT rental_id FROM rental
			WHERE return_date IS NULL
			ORDER BY random()
			LIMIT 1
		)`)
	if err != nil {
		return fmt.Errorf("update return: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil // no open rentals right now, not an error
	}
	return nil
}
