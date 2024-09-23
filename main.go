package main

import (
	"log"
	"time"

	"github.com/muhammadfarhankt/omniFlix/api"
	"github.com/muhammadfarhankt/omniFlix/db"
	"github.com/muhammadfarhankt/omniFlix/indexer"
)

func main() {
	// Initialize database connection (using db.NewDB)
	dbInstance, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbInstance.Close()

	// Create the 'blocks' table (using dbInstance.CreateTable)
	err = dbInstance.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	// Create an instance of the indexer
	idx := indexer.NewIndexer(dbInstance.DB)

	// Initialize API
	apiInstance := api.NewAPI(idx)

	// Continuous indexing (using goroutines and concurrency)
	//6341001
	minBlockHeight := int64(6341001)

	// Find already existing min height from DB

	// Find max height from REST API
	maxBlockHeight := int64(11553690)

	// Fetch the latest block height from the REST API
	latestHeight, err := idx.GetLatestBlockHeightFromREST()
	if err != nil {
		log.Printf("Error fetching latest block height: %v", err)
		latestHeight = maxBlockHeight // Example: Use the initial maxBlockHeight as a fallback
	}

	// If latestHeight is greater than maxBlockHeight, update maxBlockHeight
	if latestHeight > maxBlockHeight {
		maxBlockHeight = latestHeight
	}

	go func() {
		// infinite loop
		for {
			idx.StartIndexing(minBlockHeight, maxBlockHeight)
			time.Sleep(2 * time.Second) // Wait for 2 seconds before the next indexing cycle
		}
	}()

	// Start the API
	apiInstance.Start(":8080")
}
