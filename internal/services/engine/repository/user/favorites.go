package user

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/georgysavva/scany/pgxscan"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) AddToFavorites(ctx context.Context, data model.FavoritesData) error {
	favorites, err := r.dbGetFavorites(ctx, data.UserID)
	if err != nil && !pgxscan.NotFound(err) {
		return err
	}

	switch data.OperationType {
	case model.OperationTypePayment:
		if favorites.Payment == nil {
			favorites.Payment = make(map[string][]string)
		}

		favorites.Payment[data.ExternalSystem] = append(favorites.Payment[data.ExternalSystem], data.ExternalMethod)
	case model.OperationTypePayout:
		if favorites.Payout == nil {
			favorites.Payout = make(map[string][]string)
		}
		favorites.Payout[data.ExternalSystem] = append(favorites.Payout[data.ExternalSystem], data.ExternalMethod)
	default:
		return fmt.Errorf("unresolved operation type: %q", data.OperationType)
	}

	return r.dbSetFavorites(ctx, data.UserID, favorites)
}

func (r *Repository) RemoveFromFavorites(ctx context.Context, data model.FavoritesData) error {
	favorites, err := r.dbGetFavorites(ctx, data.UserID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil
		}
		return err
	}

	switch data.OperationType {
	case model.OperationTypePayment:
		if favorites.Payment == nil {
			return nil
		}

		favoriteExternalMethods, found := favorites.Payment[data.ExternalSystem]
		if !found {
			return nil
		}

		targetIdx := slices.Index(favoriteExternalMethods, data.ExternalMethod)
		if targetIdx == -1 {
			return nil
		}

		favorites.Payment[data.ExternalSystem] = slices.Delete(favoriteExternalMethods, targetIdx, targetIdx+1)
	case model.OperationTypePayout:
		if favorites.Payout == nil {
			return nil
		}

		favoriteExternalMethods, found := favorites.Payout[data.ExternalSystem]
		if !found {
			return nil
		}

		targetIdx := slices.Index(favoriteExternalMethods, data.ExternalMethod)
		if targetIdx == -1 {
			return nil
		}

		favorites.Payout[data.ExternalSystem] = slices.Delete(favoriteExternalMethods, targetIdx, targetIdx+1)
	default:
		return fmt.Errorf("unresolved operation type: %q", data.OperationType)
	}

	return r.dbSetFavorites(ctx, data.UserID, favorites)
}

func (r *Repository) GetFavorites(ctx context.Context, userID int64) (model.UserFavorites, error) {
	favorites, err := r.dbGetFavorites(ctx, userID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return model.UserFavorites{}, nil
		}
		return model.UserFavorites{}, err
	}
	return userFavoritesFromDB(favorites), nil
}

func (r *Repository) dbGetFavorites(ctx context.Context, userID int64) (dbUserFavorites, error) {
	var favoritesRaw []byte
	if err := pgxscan.Get(ctx, r.conn, &favoritesRaw,
		`SELECT favorites FROM "user" WHERE id = $1`,
		userID,
	); err != nil {
		return dbUserFavorites{}, err
	}

	var result dbUserFavorites
	if len(favoritesRaw) > 0 {
		if err := json.Unmarshal(favoritesRaw, &result); err != nil {
			return dbUserFavorites{}, err
		}
	}

	return result, nil
}

func (r *Repository) dbSetFavorites(ctx context.Context, userID int64, favorites dbUserFavorites) error {
	favoritesRaw, err := json.Marshal(favorites)
	if err != nil {
		return fmt.Errorf("marshal favorites %v: %w", favorites, err)
	}

	_, err = r.conn.Exec(ctx,
		`INSERT INTO "user" (id, favorites) VALUES($1, $2) ON CONFLICT (id) DO UPDATE SET favorites = $2`,
		userID, favoritesRaw,
	)
	return err
}
