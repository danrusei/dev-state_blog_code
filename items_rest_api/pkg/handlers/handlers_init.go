package handlers

import (
	"log"

	"github.com/Danr17/items-rest-api/pkg/storage"
)

//Handlers host the dependencies
type Handlers struct {
	logger    *log.Logger
	db        storage.Storage
	htmlFiles []string
}

//NewHandlers creates new Handlers struct
func NewHandlers(logger *log.Logger, db storage.Storage, files []string) *Handlers {
	return &Handlers{
		logger:    logger,
		db:        db,
		htmlFiles: files,
	}
}
