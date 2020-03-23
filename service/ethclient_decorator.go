package service

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"sync"
)

type InMemoryCacheEthClientDecorator struct {
	ethClient EthereumClient
	cache     sync.Map
}

func NewCacheEthClientDecorator(ethClient EthereumClient) *InMemoryCacheEthClientDecorator {
	log.Println("Instantiating Cached EthClient Decorator")
	return &InMemoryCacheEthClientDecorator{ethClient: ethClient}
}

func (me *InMemoryCacheEthClientDecorator) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return me.HeaderByNumber(ctx, number)
}

func (me *InMemoryCacheEthClientDecorator) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return me.BalanceAt(ctx, account, blockNumber)
}

func (me *InMemoryCacheEthClientDecorator) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	cacheKey := me.cacheKey(q)
	if cachedLogs, exists := me.cache.Load(cacheKey); exists {
		return cachedLogs.([]types.Log), nil
	}

	logs, err := me.ethClient.FilterLogs(ctx, q)
	if err != nil {
		return logs, err
	}

	me.cache.Store(cacheKey, logs)
	return logs, err
}

func (me *InMemoryCacheEthClientDecorator) cacheKey(q ethereum.FilterQuery) string {
	addressesKey := ""
	if len(q.Addresses) == 0 {
		addressesKey = ""
	}
	for _, key := range q.Addresses {
		addressesKey += "" + key.String()
	}
	return q.FromBlock.String() + "-" + q.ToBlock.String() + "-" + addressesKey
}
