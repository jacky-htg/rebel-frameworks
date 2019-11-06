package main

import (
	"context"
	"essentials/libraries/config"
	"essentials/libraries/database"
	"essentials/routing"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if _, ok := os.LookupEnv("APP_ENV"); !ok {
		config.Setup(".env")
	}

	// =========================================================================
	// Logging
	log := log.New(os.Stdout, "Essentials : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {
	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// =========================================================================

	// Start Database

	db, err := database.Open()
	if err != nil {
		return fmt.Errorf("connecting to db: %v", err)
	}
	defer db.Close()

	// parameter server
	server := http.Server{
		Addr:         os.Getenv("APP_PORT"),
		Handler:      routing.API(db, log),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	// mulai listening server
	go func() {
		log.Println("server listening on", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Membuat channel untuk mendengarkan sinyal interupsi/terminate dari OS.
	// Menggunakan channel buffered karena paket signal membutuhkannya.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Mengontrol penerimaan data dari channel,
	// jika ada error saat listenAndServe server maupun ada sinyal shutdown yang diterima
	select {
	case err := <-serverErrors:
		return fmt.Errorf("Starting server: %v", err)

	case <-shutdown:
		log.Println("caught signal, shutting down")

		// Jika ada shutdown, meminta tambahan waktu 5 detik untuk menyelesaikan proses yang sedang berjalan.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			if err := server.Close(); err != nil {
				return fmt.Errorf("could not stop server gracefully: %v", err)
			}
		}
	}

	return nil
}
