package repository

type SubscriptionRepo interface {
	StoreSubscription(address string) error
	IsSubscribed(address string) bool
}
