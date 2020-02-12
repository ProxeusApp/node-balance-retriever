package service

import (
	"context"
	"sync"
)

type (
	ethBalanceStub struct {
	}
)

func (me *ethBalanceStub) GetBalancesForAddress(ctx context.Context, _ string) (*sync.Map, error) {
	var (
		returnMap sync.Map
		returnErr error
	)

	if ctx.Value("returnMap") != nil {
		returnMap = ctx.Value("returnMap").(sync.Map)
	}

	if ctx.Value("returnErr") != nil {
		returnErr = ctx.Value("returnErr").(error)
	}

	return &returnMap, returnErr
}
