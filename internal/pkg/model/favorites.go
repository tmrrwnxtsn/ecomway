package model

type UserFavorites struct {
	Payment map[string][]string
	Payout  map[string][]string
}

type FavoritesData struct {
	OperationType  OperationType
	Currency       string
	ExternalSystem string
	ExternalMethod string
	UserID         string
}
