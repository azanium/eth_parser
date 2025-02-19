package repo

import (
	"sync"
	"testing"
)

func TestMemorySubscriptionRepo(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "store and check new address",
			address: "0x123",
			want:    true,
		},
		{
			name:    "check non-existent address",
			address: "0x456",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMemoryTransactionRepo()

			// Test StoreSubscription
			if tt.want {
				err := repo.StoreSubscription(tt.address)
				if err != nil {
					t.Errorf("StoreSubscription() error = %v", err)
				}
			}

			// Test IsSubscribed
			if got := repo.IsSubscribed(tt.address); got != tt.want {
				t.Errorf("IsSubscribed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemorySubscriptionRepoConcurrent(t *testing.T) {
	repo := NewMemoryTransactionRepo()
	addresses := []string{"0x123", "0x456", "0x789", "0xabc"}
	workers := 10

	// Test concurrent writes
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, addr := range addresses {
				err := repo.StoreSubscription(addr)
				if err != nil {
					t.Errorf("StoreSubscription() error = %v", err)
				}
			}
		}()
	}

	// Test concurrent reads while writing
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, addr := range addresses {
				_ = repo.IsSubscribed(addr)
			}
		}()
	}

	wg.Wait()

	// Verify all addresses were stored
	for _, addr := range addresses {
		if !repo.IsSubscribed(addr) {
			t.Errorf("Address %s should be subscribed", addr)
		}
	}
}
