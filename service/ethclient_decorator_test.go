package service

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCacheKey(t *testing.T) {

	ethClient := NewEthClientStub()
	ethClientDecorator := NewCacheEthClientDecorator(ethClient)

	t.Run("key is generated even with all to nil", func(t *testing.T) {
		cacheKey := ethClientDecorator.cacheKey(ethereum.FilterQuery{
			BlockHash: nil,
			FromBlock: nil,
			ToBlock:   nil,
			Addresses: nil,
			Topics:    nil,
		})

		assert.Equal(t, "<nil>-<nil>-", cacheKey)
	})

	t.Run("key with only blocks, no addresses", func(t *testing.T) {
		cacheKey := ethClientDecorator.cacheKey(ethereum.FilterQuery{
			FromBlock: big.NewInt(500),
			ToBlock:   big.NewInt(900000),
			Addresses: nil,
			Topics:    nil,
		})

		assert.Equal(t, "500-900000-", cacheKey)
	})

	t.Run("key with only blocks, no addresses", func(t *testing.T) {
		cacheKey := ethClientDecorator.cacheKey(ethereum.FilterQuery{
			FromBlock: big.NewInt(500),
			ToBlock:   big.NewInt(900000),
			Addresses: []common.Address{
				common.HexToAddress("0x043129ab3945d2bb75f3b5de21487343efbeffd2"),
				common.HexToAddress("0x123129ab3945d2bb75f3b5de21487343efbeffd2"),
			},
			Topics: nil,
		})

		assert.Equal(t, "500-900000-0x043129ab3945d2bb75f3b5de21487343efbeffd2-0x123129ab3945d2bb75f3b5de21487343efbeffd2", cacheKey)
	})

	t.Run("result is cached after first call", func(t *testing.T) {
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(500),
			ToBlock:   big.NewInt(900000),
			Addresses: []common.Address{
				common.HexToAddress("0x043129ab3945d2bb75f3b5de21487343efbeffd2"),
			},
			Topics: nil,
		}

		ethClientDecorator.FilterLogs(context.Background(), query)
		ethClientDecorator.FilterLogs(context.Background(), query)

		assert.Equal(t, 1, ethClient.FilterLogsCallsCount, "should only call FilterLogs once")
	})

}