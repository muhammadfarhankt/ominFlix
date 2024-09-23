package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarhankt/omniFlix/indexer"
)

// API struct to hold dependencies
type API struct {
	indexer *indexer.Indexer
}

// NewAPI creates a new API instance
func NewAPI(indexer *indexer.Indexer) *API {
	return &API{indexer: indexer}
}

// Start starts the API server
func (a *API) Start(addr string) {
	router := gin.Default()

	// API endpoint to fetch, compare, store, and show block details
	router.GET("/block/:height", a.getBlockDetailsHandler)

	log.Printf("Starting API server on %s", addr)
	router.Run(addr)
}

// getBlockDetailsHandler handles the /block/:height endpoint
func (a *API) getBlockDetailsHandler(c *gin.Context) {
	heightStr := c.Param("height")
	height, err := strconv.ParseInt(heightStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block height"})
		return
	}

	// Fetch block details (from DB or blockchain)
	blockDetails, err := a.indexer.GetBlockDetails(height)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blockDetails)
}
