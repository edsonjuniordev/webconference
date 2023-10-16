package store

import (
	"github.com/edsonjuniordev/webconference/domain/models"
	"github.com/edsonjuniordev/webconference/domain/store/collection"
)

type Store struct {
	RoomCollection collection.Room
}

func NewStore() Store {
	return Store{
		RoomCollection: collection.Room{
			Rooms: make(map[string]*models.Room),
		},
	}
}
