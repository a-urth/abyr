package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/a-urth/abyr/src/service/clientapi"
)

func main() {
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	service, err := clientapi.NewService()
	if err != nil {
		log.Fatal(err)
	}

	// using router here because it just looks cleaner for me that way
	router := mux.NewRouter()
	router.HandleFunc("/ping", service.Ping).Methods("GET")
	router.HandleFunc("/port/{id}", service.GetPort).Methods("GET")

	host := ":8000"
	srv := &http.Server{
		Handler:      router,
		Addr:         host,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Debugf("Starting to serve client api on %s", host)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Debug("Shutting down the server...")

	service.Close()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	srv.Shutdown(ctx)

	log.Debug("Server gracefully stopped")
}
