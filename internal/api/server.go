package api

import (
	"sync/atomic"

	"github.com/memartello/go-http-server/internal/database"
)

type API struct {
	dbQueries *database.Queries
	fileServerHits atomic.Int32
}

func NewAPI(dbQueries *database.Queries) *API {	
	return &API{dbQueries: dbQueries, fileServerHits: atomic.Int32{}}
}

// Initialize Server and handle routes here.