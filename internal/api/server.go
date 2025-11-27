package api

import (
	"sync/atomic"

	"github.com/memartello/go-http-server/internal/database"
)

type API struct {
	dbQueries *database.Queries
	fileServerHits atomic.Int32
	secret string
}

func NewAPI(dbQueries *database.Queries, secret string) *API {	
	return &API{dbQueries: dbQueries, fileServerHits: atomic.Int32{}, secret: secret}
}

// Initialize Server and handle routes here.