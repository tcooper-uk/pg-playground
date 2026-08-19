package simulator

import (
	"math/rand"
	"sync"

	"github.com/tcooper/pg-playground/app/config"
)

type WeightStore struct {
	mu      sync.RWMutex
	weights config.WeightConfig
}

func NewWeightStore(initial config.WeightConfig) *WeightStore {
	return &WeightStore{weights: initial}
}

func (w *WeightStore) Get() config.WeightConfig {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.weights
}

func (w *WeightStore) Set(wc config.WeightConfig) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.weights = wc
}

// Pick samples an op type proportionally from the current weights.
func (w *WeightStore) Pick() opType {
	w.mu.RLock()
	wc := w.weights
	w.mu.RUnlock()

	total := wc.Rental + wc.Return + wc.CustomerChurn + wc.Read
	if total == 0 {
		return opRead
	}

	n := rand.Intn(total)
	if n < wc.Rental {
		return opRental
	}
	n -= wc.Rental
	if n < wc.Return {
		return opReturn
	}
	n -= wc.Return
	if n < wc.CustomerChurn {
		return opChurn
	}
	return opRead
}
