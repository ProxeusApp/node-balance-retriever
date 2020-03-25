package service

import (
	"context"
	"sync"
)

type EthBalanceService interface {
	GetBalancesForAddress(ctx context.Context, address string) (*sync.Map, error)
}
