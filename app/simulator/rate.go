package simulator

import "sync/atomic"

type RateStore struct {
	rps atomic.Int64
}

func NewRateStore(initial int) *RateStore {
	r := &RateStore{}
	r.rps.Store(int64(initial))
	return r
}

func (r *RateStore) Get() int { return int(r.rps.Load()) }
func (r *RateStore) Set(rps int) {
	if rps < 1 {
		rps = 1
	}
	r.rps.Store(int64(rps))
}
