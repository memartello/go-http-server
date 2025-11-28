package api

import (
	"sync/atomic"

	"github.com/memartello/go-http-server/internal/database"
)

type API struct {
	dbQueries *database.Queries
	fileServerHits atomic.Int32
	secret string
	polka_key string
}

func NewAPI(dbQueries *database.Queries, secret, polka_key string) *API {	
	return &API{dbQueries: dbQueries, fileServerHits: atomic.Int32{}, secret: secret, polka_key: polka_key}
}

// Initialize Server and handle routes here.