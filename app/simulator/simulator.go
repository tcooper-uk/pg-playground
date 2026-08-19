package simulator

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tcooper/pg-playground/app/config"
)

type Metrics struct {
	Rental atomic.Int64
	Return atomic.Int64
	Churn  atomic.Int64
	Read   atomic.Int64
	Errors atomic.Int64
}

type MetricsSnapshot struct {
	Rental int64 `json:"rental"`
	Return int64 `json:"return"`
	Churn  int64 `json:"churn"`
	Read   int64 `json:"read"`
	Errors int64 `json:"errors"`
}

func NewMetrics() *Metrics { return &Metrics{} }

func (m *Metrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		Rental: m.Rental.Load(),
		Return: m.Return.Load(),
		Churn:  m.Churn.Load(),
		Read:   m.Read.Load(),
		Errors: m.Errors.Load(),
	}
}

type lookupIDs struct {
	mu          sync.RWMutex
	customerIDs []int
	inventoryIDs []int
	staffIDs    []int
	storeIDs    []int
	cityIDs     []int
	filmIDs     []int
	actorIDs    []int
}

func (l *lookupIDs) randomCustomerID() int  { return l.pick(l.customerIDs) }
func (l *lookupIDs) randomInventoryID() int { return l.pick(l.inventoryIDs) }
func (l *lookupIDs) randomStaffID() int     { return l.pick(l.staffIDs) }
func (l *lookupIDs) randomStoreID() int     { return l.pick(l.storeIDs) }
func (l *lookupIDs) randomCityID() int      { return l.pick(l.cityIDs) }
func (l *lookupIDs) randomFilmID() int      { return l.pick(l.filmIDs) }
func (l *lookupIDs) randomActorID() int     { return l.pick(l.actorIDs) }

func (l *lookupIDs) pick(ids []int) int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return ids[rand.Intn(len(ids))]
}

type opType int

const (
	opRental opType = iota
	opReturn
	opChurn
	opRead
)

func Run(ctx context.Context, cfg config.SimulatorConfig, primary, replica *pgxpool.Pool, m *Metrics, ws *WeightStore, rs *RateStore) {
	ids, err := loadLookupIDs(ctx, primary)
	if err != nil {
		log.Printf("simulator: failed to load lookup IDs: %v", err)
		return
	}

	var wg sync.WaitGroup
	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				op := ws.Pick()
				var err error
				switch op {
				case opRental:
					err = simulateRental(ctx, primary, ids)
					if err == nil { m.Rental.Add(1) }
				case opReturn:
					err = simulateReturn(ctx, primary)
					if err == nil { m.Return.Add(1) }
				case opChurn:
					err = simulateChurn(ctx, primary, ids)
					if err == nil { m.Churn.Add(1) }
				case opRead:
					err = simulateRead(ctx, replica, ids)
					if err == nil { m.Read.Add(1) }
				}
				if err != nil {
					m.Errors.Add(1)
					log.Printf("simulator op %d: %v", op, err)
				}

				sleep := time.Duration(cfg.Workers) * time.Second / time.Duration(rs.Get())
				select {
				case <-ctx.Done():
					return
				case <-time.After(sleep):
				}
			}
		}()
	}

	wg.Wait()
}


func loadLookupIDs(ctx context.Context, pool *pgxpool.Pool) (*lookupIDs, error) {
	ids := &lookupIDs{}

	queries := []struct {
		query string
		dest  *[]int
	}{
		{`SELECT customer_id FROM customer WHERE activebool = true`, &ids.customerIDs},
		{`SELECT inventory_id FROM inventory`, &ids.inventoryIDs},
		{`SELECT staff_id FROM staff`, &ids.staffIDs},
		{`SELECT store_id FROM store`, &ids.storeIDs},
		{`SELECT city_id FROM city`, &ids.cityIDs},
		{`SELECT film_id FROM film`, &ids.filmIDs},
		{`SELECT actor_id FROM actor`, &ids.actorIDs},
	}

	for _, q := range queries {
		rows, err := pool.Query(ctx, q.query)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return nil, err
			}
			*q.dest = append(*q.dest, id)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return ids, nil
}
