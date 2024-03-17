package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cody-Kao/crud-sql/db"
	"github.com/Cody-Kao/crud-sql/handlers"
	"github.com/Cody-Kao/crud-sql/server"
)

func main() {
	// Connect to the database
	DB, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	fmt.Println("Db connection established")
	// Set up signal handling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Set up signal handler
	go func() {
		<-interrupt
		fmt.Println("Shutting down...")
		handlers.QueryOne.Close()
		handlers.QueryAll.Close()
		handlers.CreateRow.Close()
		handlers.UpdateRow.Close()
		handlers.DeleteRow.Close()
		fmt.Println("All query statements are closed")

		err := DB.Close()
		if err != nil {
			log.Fatal("Error when closing DB", err)
		}
		fmt.Println("DB closed gracefully")
		cancel()
	}()

	// Pass the DB variable to the handler package
	handlers.SetDB(DB)

	// Create server
	srv := server.CreateServer()

	// Start the server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for the shutdown signal and block the above go routine for preventing early stop
	<-ctx.Done()

	// Shutdown the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
