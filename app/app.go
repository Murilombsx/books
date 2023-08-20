package app

import (
	"books/api"
	"books/constants"
	"books/dataprovider"
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	server *http.Server
}

func (a *App) Init() {
	router := mux.NewRouter()
	router.Handle("/books", api.NewBooksHandler(dataprovider.NewDataProvider(constants.URL, &http.Client{Timeout: constants.DEFAULT_TIMEOUT}))).Methods(http.MethodGet)
	router.Handle("/health", api.NewHeatlhHandler()).Methods(http.MethodGet)
	router.Use(loggerMiddleware)

	a.server = &http.Server{
		Handler:      router,
		Addr:         constants.SERVER_ADDR,
		ReadTimeout:  constants.DEFAULT_TIMEOUT,
		WriteTimeout: constants.DEFAULT_TIMEOUT,
	}
}

func (a *App) Start() {
	go func() {
		log.Info("starting server")
		if err := a.server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

func (a *App) GracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()
	a.server.Shutdown(ctx)

	log.Info("shutting down")
	os.Exit(0)
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"path":       r.URL.Path,
			"parameters": r.URL.RawQuery,
		}).Info("starting request")

		startTime := time.Now()
		next.ServeHTTP(w, r)

		log.WithFields(log.Fields{
			"path":       r.URL.Path,
			"parameters": r.URL.RawQuery,
			"duration":   time.Since(startTime).String(),
		}).Info("completing request")
	})
}
