package prices

import (
	"fmt"
	"sync"
	"time"

	crrqueue "github.com/Workiva/go-datastructures/queue"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

const (
	_DEFAULT_CAPACITY = 256
	_SLEEP_INTERVAL   = 5000
)

var (
	priceBuffer *crrqueue.Queue
	priceLock   sync.Mutex
)

func Initialize() {
	if priceBuffer != nil {
		logrus.Warn(logger.PRICES_DOUBLE_INITIALIZATION)
		return
	}
	priceBuffer = crrqueue.New(_DEFAULT_CAPACITY)

	go func() {
		ticker := time.NewTicker(_SLEEP_INTERVAL * time.Millisecond)
		for range ticker.C {
			if store(false) < 0 {
				return
			}
		}
	}()
}

func InsertMany(prices []model.SymbolPrice) error {
	return insert_many(prices)
}

func InsertManyDeferred(prices []model.SymbolPrice) error {
	if priceBuffer == nil {
		err := fmt.Errorf(logger.ANAL_PRICES_ERR_NO_INITIALIZATION)
		logrus.Error(err.Error())
		return err
	}

	payload := make([]interface{}, 0, len(prices))
	for _, price := range prices {
		payload = append(payload, price)
	}

	err := priceBuffer.Put(payload...)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func Get(symbols []string, start, end int64) ([]model.SymbolPrice, error) {
	return find(symbols, start, end)
}

func GetByTimestamp(symbols []string, start, end int64) ([]model.SymbolPriceByTimestamp, error) {
	return find_by_timestamp(symbols, start, end)
}

func Terminate() {
	if priceBuffer == nil {
		logrus.Warn(logger.PRICES_NO_INITIALIZATION)
		return
	}
	store(true)
}

func store(stop bool) int {
	// If queue was disposed, return -1
	if priceBuffer.Disposed() {
		return -1
	}

	// Protecting critical sequence of code
	priceLock.Lock()
	defer priceLock.Unlock()

	// Retrieving items to store
	var err error
	var iitems []interface{}
	if stop {
		iitems = priceBuffer.Dispose()
	} else {
		iitems, err = priceBuffer.Get(priceBuffer.Len())
	}

	// Check error
	if err != nil {
		logrus.Error(err.Error())
	}

	// If nothing to store, return
	if len(iitems) == 0 {
		return 0
	}

	// Store items
	items := make([]model.SymbolPrice, 0, len(iitems))
	for _, iitem := range iitems {
		items = append(items, iitem.(model.SymbolPrice))
	}
	logrus.Infof(logger.PRICES_INSERT_MANY, len(items))
	insert_many(items)
	return len(items)
}
