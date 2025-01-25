package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)

	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	cfg.dbQueries.Reset(context.Background())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
