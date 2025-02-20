package engine

import (
	"context"
	"multitenant/model"
	"multitenant/pkg/cache"
	"multitenant/pkg/config"
	"multitenant/pkg/db"
	"multitenant/pkg/logger"
	"multitenant/pkg/metrics"
	"sync"
	"time"
)

const (
	batchSize              = 100             // Adjust for optimal performance
	flushInterval          = 2 * time.Second // Flush batch every X seconds
	workerPoolSize         = 10              // Set the number of worker goroutines
	autoRefreshTxThreshold = 20              // Refresh cache after X transactions
)

type Engine struct {
	config      *config.Config
	logger      *logger.Logger
	db          *db.DB
	cache       *cache.Cache
	metrics     *metrics.Metrics
	txQueue     chan *model.Transaction
	workerQueue chan []*model.Transaction
	batchBuffer []*model.Transaction
	mu          sync.Mutex
	wg          sync.WaitGroup
	mp          sync.Map // for caching tenant and product existence
	ctx         context.Context
}

func NewEngine(ctx context.Context, cfg *config.Config, l *logger.Logger, db *db.DB, c *cache.Cache, m *metrics.Metrics) *Engine {

	return &Engine{
		ctx:         ctx,
		config:      cfg,
		logger:      l,
		db:          db,
		cache:       c,
		metrics:     m,
		txQueue:     make(chan *model.Transaction, batchSize),
		batchBuffer: make([]*model.Transaction, 0, batchSize),
		workerQueue: make(chan []*model.Transaction, workerPoolSize),
		mu:          sync.Mutex{},
		mp:          sync.Map{},
		wg:          sync.WaitGroup{},
	}
}

func (e *Engine) Start() error {
	e.logger.Debugf("Starting engine")
	e.startWorkerPool(workerPoolSize)
	e.wg.Add(1)
	go e.processBatches()
	return nil
}

func (e *Engine) Stop() {
	e.logger.Debugf("Stopping engine")
	e.wg.Wait() // block until all workers are done, after context is cancelled
	close(e.txQueue)
}
