package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"

	"github.com/danrusei/items-rest-api/pkg/handlers"
	"github.com/danrusei/items-rest-api/pkg/storage"
	"github.com/danrusei/items-rest-api/pkg/storage/dbfirestore"
)

var (
	listenAddr string
	dbType     string
)

// api holds dependencies
type api struct {
	mutex     sync.Mutex
	db        storage.Storage
	router    *http.ServeMux
	logger    *log.Logger
	htmlFiles []string
}

func newAPI() *api {
	a := &api{
		router: http.NewServeMux(),
		mutex:  sync.Mutex{},
	}
	return a
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.StringVar(&listenAddr, "listen-addr", ":8080", "server listen address")
	flag.StringVar(&dbType, "db-type", "firestore", "select database, memory or firestore")
	flag.Parse()

	API := newAPI()

	API.logger = log.New(os.Stdout, "gcuk ", log.LstdFlags|log.Lshortfile)
	API.htmlFiles = []string{"../tmpl/index.html"}

	//add firestore
	ctx := context.Background()
	var firestoreClient *firestore.Client
	// get_service_account_key.json it is just a fake key, you need a real one
	sa := option.WithCredentialsFile("../key/get_service_account_key.json")
	// Your_GCP_Account -- use your real GCP Account ID
	firestoreClient, err := firestore.NewClient(ctx, "Your_GCP_Account", sa)
	if err != nil {
		log.Printf("Client initialization error: %v", err)
		os.Exit(1)
	}

	switch dbType {
	//	case "memory":
	//		API.db = dbmemory.NewMemory()
	case "firestore":
		API.db = dbfirestore.NewFirestoreDB(firestoreClient)
	default:
		API.db = dbfirestore.NewFirestoreDB(firestoreClient)
	}

	h := handlers.NewHandlers(API.logger, API.db, API.htmlFiles)

	mux := API.router
	h.CreateRoutes(mux)

	//storage.PopulateItems(API.db)

	server := http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main : API listening on %s", listenAddr)
		serverErrors <- server.ListenAndServe()
	}()

	//channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	//blocking run and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting server: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown")

		//give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// asking listener to shutdown
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = server.Close()
		}

		if err != nil {
			return fmt.Errorf("main : could not stop server gracefully : %v", err)
		}
	}

	return nil
}
