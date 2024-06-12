package user

import (
	"maps"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type dbUserFavorites struct {
	Payment map[string][]string `json:"payment"`
	Payout  map[string][]string `json:"payout"`
}

func userFavoritesFromDB(dbFavorites dbUserFavorites) model.UserFavorites {
	return model.UserFavorites{
		Payment: maps.Clone(dbFavorites.Payment),
		Payout:  maps.Clone(dbFavorites.Payout),
	}
}
