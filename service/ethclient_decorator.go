package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"io"
	"log"
	"math/big"
	"os"
	"sync"
)

type diskCacheEthClientDecorator struct {
	ethClient EthereumClient
	cache     sync.Map
	diskLock  sync.Mutex

	filename string
}

func NewDiskCacheEthClientDecorator(ethClient EthereumClient, filename string) (*diskCacheEthClientDecorator, error) {
	inMemoryCacheEthClientDecorator := diskCacheEthClientDecorator{ethClient: ethClient, filename: filename}
	err := inMemoryCacheEthClientDecorator.restoreCacheFromFile()
	if err != nil {
		return nil, err
	}

	return &inMemoryCacheEthClientDecorator, nil
}

func (me *diskCacheEthClientDecorator) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return me.ethClient.HeaderByNumber(ctx, number)
}

func (me *diskCacheEthClientDecorator) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return me.ethClient.BalanceAt(ctx, account, blockNumber)
}

func (me *diskCacheEthClientDecorator) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	cacheKey := me.cacheKey(q)
	if cachedLogs, exists := me.cache.Load(cacheKey); exists {
		log.Println("Cached, return data for cache key " + cacheKey)
		return cachedLogs.([]types.Log), nil
	}

	logs, err := me.ethClient.FilterLogs(ctx, q)
	if err != nil {
		return logs, err
	}

	log.Println("Not cached, store it for cache key " + cacheKey)
	me.cache.Store(cacheKey, logs)
	return logs, err
}

func (me *diskCacheEthClientDecorator) cacheKey(q ethereum.FilterQuery) string {
	addressesKey := ""
	if len(q.Addresses) == 0 {
		addressesKey = ""
	}
	for _, key := range q.Addresses {
		addressesKey += "" + key.String()
	}
	return q.FromBlock.String() + "-" + q.ToBlock.String() + "-" + addressesKey
}

func (me *diskCacheEthClientDecorator) restoreCacheFromFile() error {
	me.diskLock.Lock()
	defer me.diskLock.Unlock()

	f, err := os.Open(me.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(me.cache)
	if err != nil {
		return err
	}

	return nil
}

func (me *diskCacheEthClientDecorator) updateCacheFile() error {
	me.diskLock.Lock()
	defer me.diskLock.Unlock()

	f, err := os.Create(me.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := json.Marshal(me.cache)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(data))
	return err
}
