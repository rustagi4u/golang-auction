package main

import (
	"context"
	"controller"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shared"
	database "shared"
	"time"

	"github.com/gorilla/mux"
)

// DBInitializerCaller starts Postgres db
func DBInitializerCaller() {
	DB, err := database.NewOpen()

	if err != nil {
		panic(err.Error())
	}

	defer DB.Close()

	// Open doesn’t open a connection. Validate DSN data:
	err = DB.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	fmt.Println("Connected to database")

}

// Server starts the app on a port
func main() {
	var wait time.Duration
	DBInitializerCaller()

	router := mux.NewRouter()

	// router.HandleFunc("/clients", controller.GetClients).Methods("GET")
	// router.HandleFunc("/client", controller.CreateClient).Methods("POST")
	router.HandleFunc("/api/v1/offer", controller.CreateOffer).Methods("POST")
	router.HandleFunc("/api/v1/offer", controller.GetOffer).Methods("Get")
	router.HandleFunc("/api/v1/bid", controller.CreateBid).Methods("POST")
	router.HandleFunc("/api/user/login", controller.Authenticate).Methods("POST")
	router.HandleFunc("/api/user/new", controller.CreateAccount).Methods("POST")
	router.HandleFunc("/api/v1/sold", controller.SoldOffers).Methods("GET")
	router.HandleFunc("/api/v1/accepted", controller.BidAccepted).Methods("PUT")
	router.HandleFunc("/", controller.RootHandler).Methods("GET")
	router.HandleFunc("/ws", controller.WsHandler)
	go controller.Echo()

	router.Use(shared.JwtAuthentication)
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:7000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  120 * time.Minute,
	}

	// Run our server in a goroutine so that it doesn’t block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn’t block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
