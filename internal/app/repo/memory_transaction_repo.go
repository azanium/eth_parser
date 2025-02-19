package repo

import (
	"eth_parser/internal/domain/repository"
	"sync"
)

var _ repository.SubscriptionRepo = (*MemorySubscriptionRepo)(nil)

type MemorySubscriptionRepo struct {
	subscriptions map[string]bool
	mutex         sync.RWMutex
}

func NewMemoryTransactionRepo() *MemorySubscriptionRepo {
	return &MemorySubscriptionRepo{
		subscriptions: make(map[string]bool),
	}
}

func (r *MemorySubscriptionRepo) StoreSubscription(address string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.subscriptions[address] = true
	return nil
}

func (r *MemorySubscriptionRepo) IsSubscribed(address string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, ok := r.subscriptions[address]
	return ok
}
