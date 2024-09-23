package indexer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// BlockDetails represents the structure for block data
type BlockDetails struct {
	Height          int64           `json:"height"`
	BlockID         string          `json:"block_id"`
	NumTransactions int             `json:"num_transactions"`
	Proposer        string          `json:"proposer"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       sql.NullTime    `json:"deleted_at"`
	Details         json.RawMessage `json:"details"`
}

// Indexer struct to hold dependencies
type Indexer struct {
	db *sql.DB
}

// NewIndexer creates a new Indexer instance
func NewIndexer(db *sql.DB) *Indexer {
	return &Indexer{
		db: db,
	}
}

// GetBlockDetails fetches block details from the database if available,
// otherwise fetches from the blockchain and stores it in the database.
func (idx *Indexer) GetBlockDetails(height int64) (*BlockDetails, error) {
	// 1. Try fetching from Postgres first
	var blockDetails BlockDetails
	err := idx.db.QueryRow("SELECT block_height, block_id,proposer_address, num_transactions, created_at, updated_at, deleted_at, details FROM blocks WHERE block_height = $1", height).Scan(
		&blockDetails.Height,
		&blockDetails.BlockID,
		&blockDetails.Proposer,
		&blockDetails.NumTransactions,
		&blockDetails.CreatedAt,
		&blockDetails.UpdatedAt,
		&blockDetails.DeletedAt,
		&blockDetails.Details,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// 2. If not found in Postgres, fetch from blockchain
			blockDetails, err = idx.FetchAndStoreBlockDetails(height)
			if err != nil {
				return nil, fmt.Errorf("error fetching and storing block details: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error fetching block details from database: %w", err)
		}
	}

	return &blockDetails, nil
}

// StartIndexing starts the continuous indexing process with concurrency
func (idx *Indexer) StartIndexing(minBlockHeight, maxBlockHeight int64) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 100) // Limit concurrency to 100 goroutines

	// Fetch the latest block height
	latestHeight, err := idx.GetLatestBlockHeightFromREST()
	if err != nil {
		log.Printf("Error fetching latest block height: %v", err)
	}

	// If latestHeight is greater than maxBlockHeight, update maxBlockHeight
	if latestHeight > maxBlockHeight {
		log.Printf("Latest block height %d is greater than max block height %d, updating max block height to %d", latestHeight, maxBlockHeight, latestHeight)
		maxBlockHeight = latestHeight
	}

	for currentHeight := maxBlockHeight; currentHeight >= minBlockHeight; currentHeight-- {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire a semaphore slot

		go func(height int64) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release the semaphore slot

			_, err := idx.FetchAndStoreBlockDetails(height)
			if err != nil {
				log.Printf("Error indexing block %d: %v", height, err)
			}
		}(currentHeight)
	}

	wg.Wait()
}

// FetchAndStoreBlockDetails fetches and stores block details with timestamps (using only RPC)
func (idx *Indexer) FetchAndStoreBlockDetails(height int64) (BlockDetails, error) {
	var (
		blockDetails BlockDetails
		err          error
	)

	blockDetails, err = idx.getBlockResults(height)
	if err != nil {
		return BlockDetails{}, fmt.Errorf("error getting block details: %w", err)
	}

	// Store blockDetails in the database with timestamps
	go func() {
		ctx := context.Background()

		// Convert Details to JSON string
		detailsJSON, err := json.Marshal(blockDetails.Details)
		if err != nil {
			log.Printf("Error marshaling details to JSON: %v", err)
			return
		}

		currentTime := time.Now()
		_, err = idx.db.ExecContext(ctx, `
			INSERT INTO blocks (block_height, block_id, proposer_address, num_transactions, details, created_at, updated_at, deleted_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, NULL)
			ON CONFLICT (block_height) DO UPDATE 
			SET block_id = EXCLUDED.block_id,
				proposer_address = EXCLUDED.proposer_address,
				num_transactions = EXCLUDED.num_transactions,
				details = EXCLUDED.details,
				updated_at = EXCLUDED.updated_at`,
			height, blockDetails.BlockID, blockDetails.Proposer, blockDetails.NumTransactions, detailsJSON, currentTime, currentTime)
		if err != nil {
			log.Printf("Error storing block data in database: %v", err)
		}
	}()

	return blockDetails, nil
}

// GetLatestBlockHeight fetches the latest block height
func (idx *Indexer) GetLatestBlockHeight() (int64, error) {
	url := "https://rpc.omniflix.network/status"
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error fetching status: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding status: %w", err)
	}

	syncInfo, ok := result["result"].(map[string]interface{})["sync_info"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("sync_info not found in response")
	}
	latestBlockHeight, ok := syncInfo["latest_block_height"].(string)
	if !ok {
		return 0, fmt.Errorf("latest_block_height not found in response")
	}
	height, err := strconv.ParseInt(latestBlockHeight, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing latest block height: %w", err)
	}

	return height, nil
}

// GetLatestBlockHeightFromREST fetches the latest block height from the REST API
func (idx *Indexer) GetLatestBlockHeightFromREST() (int64, error) {
	url := "https://rest.omniflix.network/cosmos/base/tendermint/v1beta1/blocks/latest"

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error fetching latest block from REST API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("REST API request failed with status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding latest block from REST API: %w", err)
	}

	// Extract block height (adapt based on actual JSON structure)
	block, ok := result["block"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid REST API response: 'block' field not found")
	}
	header, ok := block["header"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid REST API response: 'header' field not found")
	}
	heightStr, ok := header["height"].(string)
	if !ok {
		return 0, fmt.Errorf("invalid REST API response: 'height' field not found")
	}

	height, err := strconv.ParseInt(heightStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing block height: %w", err)
	}

	return height, nil
}

// getBlockResults fetches block results from the RPC /block_results endpoint
func (idx *Indexer) getBlockResults(height int64) (BlockDetails, error) {
	url := fmt.Sprintf("https://rpc.omniflix.network/block_results?height=%d", height)
	resp, err := http.Get(url)
	if err != nil {
		return BlockDetails{}, fmt.Errorf("error fetching block results: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return BlockDetails{}, fmt.Errorf("error decoding block results: %w", err)
	}

	resultResult, ok := result["result"].(map[string]interface{})
	if !ok || resultResult == nil {
		return BlockDetails{}, fmt.Errorf("invalid or missing 'result' field in /block_results API response")
	}

	// If block_id is not found in /block_results, try fetching it from /block
	blockData, err := idx.getBlock(height)
	if err != nil {
		return BlockDetails{}, fmt.Errorf("error fetching block_id from /block: %w", err)
	}
	// Use the block_id and proposer from /block response
	blockID := blockData.BlockID
	proposer := blockData.Proposer

	txsResults, ok := resultResult["txs_results"]
	if !ok {
		return BlockDetails{}, fmt.Errorf("error extracting txs_results from block results")
	}

	numTransactions := 0
	switch txs := txsResults.(type) {
	case []interface{}:
		numTransactions = len(txs)
	case nil:
		numTransactions = 0
	default:
		return BlockDetails{}, fmt.Errorf("unexpected type for txs_results: %T", txs)
	}

	blockDetails := BlockDetails{
		Height:          height,
		BlockID:         blockID,
		Proposer:        proposer,
		NumTransactions: numTransactions,
		// ... extract details, etc. from resultResult
	}
	return blockDetails, nil
}

// getBlock fetches block data from the RPC /block endpoint (for extracting block_id and proposer)
func (idx *Indexer) getBlock(height int64) (BlockDetails, error) {
	url := fmt.Sprintf("https://rpc.omniflix.network/block?height=%d", height)
	resp, err := http.Get(url)
	if err != nil {
		return BlockDetails{}, fmt.Errorf("error fetching block from RPC: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return BlockDetails{}, fmt.Errorf("error decoding block from RPC: %w", err)
	}

	resultResult, ok := result["result"].(map[string]interface{})
	if !ok || resultResult == nil {
		return BlockDetails{}, fmt.Errorf("invalid or missing 'result' field in /block API response")
	}

	blockID, ok := resultResult["block_id"].(map[string]interface{})["hash"].(string)
	if !ok {
		return BlockDetails{}, fmt.Errorf("error extracting block_id from /block response")
	}

	proposer, ok := resultResult["block"].(map[string]interface{})["header"].(map[string]interface{})["proposer_address"].(string)
	if !ok {
		return BlockDetails{}, fmt.Errorf("error extracting proposer_address from /block response")
	}

	blockDetails := BlockDetails{
		BlockID:  blockID,
		Proposer: proposer,
	}
	return blockDetails, nil
}
