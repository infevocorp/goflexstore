package sql

func NewStores() *Stores {
	return &Stores{}
}

type Stores struct {
	User UserStore
}
