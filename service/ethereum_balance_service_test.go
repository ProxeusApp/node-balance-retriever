package service

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
)

func TestDefaultTaxReporterService_GetBalances(t *testing.T) {
	ethBalanceStub := &ethBalanceStub{}
	taxReporter := NewEthereumBalanceService(ethBalanceStub)

	t.Run("ShouldReturnETHandXesBalance", func(t *testing.T) {

		returnMap := sync.Map{}
		eth := big.NewInt(1231230982)

		xes, ok := big.NewInt(0).SetString("278797678000000000000000000", 10)
		if !ok {
			t.Errorf("big int error")
		}
		returnMap.Store("ETH", eth)
		returnMap.Store("XES", xes)

		ctx := context.WithValue(context.Background(), "returnMap", returnMap)
		ctx = context.WithValue(ctx, "returnErr", nil)

		taxReporterBalances, err := taxReporter.GetBalances(ctx, "0x1")

		if err != nil {
			t.Error(err)
		}
		if taxReporterBalances["ETH"].Cmp(big.NewFloat(0.000000001231230982)) != 0 {
			t.Errorf("expected ETH to be %s but got %s", "0.000000001231230982", taxReporterBalances["ETH"])
		}
		if taxReporterBalances["XES"].Cmp(big.NewFloat(278797678)) != 0 {
			t.Errorf("expected XES to be %s but got %s", "278797678", taxReporterBalances["XES"])
		}
		if taxReporterBalances["MKR"] != nil {
			t.Error("Expected MKR to be nil")
		}
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		expectedError := errors.New("eth error")
		ctx := context.WithValue(context.Background(), "returnErr", expectedError)
		taxReporterBalances, err := taxReporter.GetBalances(ctx, "0x1")

		if err != expectedError {
			t.Error("Expected err but was nil")
		}
		if taxReporterBalances != nil {
			t.Error("Expected taxReporterBalances to be nil")
		}
	})
}

func TestConvertToDefaultUnit(t *testing.T) {
	taxReporter := NewEthereumBalanceService(nil)
	t.Run("ShouldConvertUnit", func(t *testing.T) {
		res := taxReporter.convertToDefaultUnit(big.NewInt(10000000000000000))
		if res.Cmp(big.NewFloat(0.01)) != 0 {
			t.Error("err 1")
		}

		if res.String() != "0.01" {
			t.Error("err 2")
		}
	})
}
