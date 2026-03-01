// Command vibium-mcp provides an MCP server for browser automation.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/plexusone/vibium-go/mcp"
)

func main() {
	headless := flag.Bool("headless", true, "Run browser in headless mode")
	project := flag.String("project", "vibium-tests", "Project name for reports")
	timeout := flag.Duration("timeout", 30*time.Second, "Default timeout for browser operations")
	flag.Parse()

	config := mcp.Config{
		Headless:       *headless,
		Project:        *project,
		DefaultTimeout: *timeout,
	}

	server := mcp.NewServer(config)

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutting down...")
		cancel()
		if err := server.Close(context.Background()); err != nil {
			log.Printf("Error closing server: %v", err)
		}
	}()

	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
