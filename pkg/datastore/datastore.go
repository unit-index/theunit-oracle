package datastore

type Datastore interface {
	Start() error
	Wait()
	Prices() PriceStore
}
