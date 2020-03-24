package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"os"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type diskCacheEthClientDecorator struct {
	ethClient EthereumClient
	cache     sync.Map
	diskLock  sync.Mutex

	filename string
}

func NewDiskCacheEthClientDecorator(ethClient EthereumClient, filename string) (*diskCacheEthClientDecorator, error) {
	client := diskCacheEthClientDecorator{ethClient: ethClient, filename: filename}
	/*if !client.fileExists() {
		_, err := os.Create(client.filename)
		if err != nil {
			log.Println("Can't create file " + client.filename)
			return nil, err
		}
	} else {
		err := client.restoreCacheFromFile()
		if err != nil {
			log.Println("Can't restore cache from file", err)
			return nil, err
		}
	}*/

	return &client, nil
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

	me.sortAddresses(&q)
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

func (me *diskCacheEthClientDecorator) fileExists() bool {
	info, err := os.Stat(me.filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
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

func (me *diskCacheEthClientDecorator) sortAddresses(query *ethereum.FilterQuery) {
	stringAddresses := make([]string, len(query.Addresses))
	for i, address := range query.Addresses {
		stringAddresses[i] = address.String()
	}

	sort.Strings(stringAddresses)

	for i, stringAddress := range stringAddresses {
		query.Addresses[i] = common.HexToAddress(stringAddress)
	}
}
