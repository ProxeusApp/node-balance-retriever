package service

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"time"
)

type (
	EthereumBalanceService interface {
		GetBalances(ctx context.Context, ethAddress string) (map[string]*big.Float, error)
	}

	defaultEthereumBalanceService struct {
		ethBalanceService EthBalanceService
	}
)

const defaultEthereumUnit = 1000000000000000000

func NewEthereumBalanceService(ethBalanceService EthBalanceService) *defaultEthereumBalanceService {
	return &defaultEthereumBalanceService{ethBalanceService: ethBalanceService}
}

// Returns the balance of tokens in a map. Are converted to default unit, see `defaultEthereumUnit`.
// This method is only compatible for erc20 tokens that use `defaultEthereumUnit` and ETH
func (me *defaultEthereumBalanceService) GetBalances(ctx context.Context, ethAddress string) (map[string]*big.Float, error) {
	response := make(map[string]*big.Float)

	ctx, cancel := context.WithTimeout(ctx, time.Minute*10)
	balances, err := me.ethBalanceService.GetBalancesForAddress(ctx, ethAddress)
	cancel()
	if err != nil {
		return nil, err
	}

	balances.Range(func(key, value interface{}) bool {
		keyString, ok := key.(string)
		if !ok {
			err = fmt.Errorf("key expected to be a string, was %v. Value: %v", reflect.TypeOf(key), key)
			return false
		}

		valWei, ok := value.(*big.Int)
		if !ok {
			err = fmt.Errorf("[taxreporter][next] casting error: %s", value)
			return false
		}

		response[keyString] = me.convertToDefaultUnit(valWei)
		return true
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (me *defaultEthereumBalanceService) convertToDefaultUnit(value *big.Int) *big.Float {
	val, ok := big.NewFloat(0).SetString(value.String())
	if !ok {
		return big.NewFloat(0)
	}
	return big.NewFloat(0).Quo(val, big.NewFloat(defaultEthereumUnit))
}
