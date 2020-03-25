package service

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBalancesForAddress(t *testing.T) {
	balanceXES := big.NewInt(0)
	balanceXES.SetString("5483000000000000000000", 10)
	json := ethplorerResponse{
		Address: "",
		ETH: ethBalance{
			Balance: 5.05,
		},
		CountTxs: 0,
		Tokens: []token{
			{
				tokenInfo: tokenInfo{
					Symbol: "XES",
				},
				Balance:  BigInt{BigIntValue: balanceXES},
				TotalIn:  0,
				TotalOut: 0,
			},
			{
				tokenInfo: tokenInfo{
					Symbol: "MKR",
				},
				Balance:  BigInt{BigIntValue: big.NewInt(7373767001504)},
				TotalIn:  0,
				TotalOut: 0,
			},
		},
	}

	tokensMap := map[string]string{
		"0x84E0b37e8f5B4B86d5d299b0B0e33686405A3919": "XES",
		"0x710129558E8ffF5caB9c0c9c43b99d79Ed864B99": "MKR",
		"0x123456558E8ffF5caB9c0c9c43b99d79Ed864B99": "ANY",
	}
	balanceService := NewEthplorerBalanceService(tokensMap)
	balances := balanceService.toMap(json)

	xesBalance, _ := balances.Load("XES")
	mkrBalance, _ := balances.Load("MKR")
	anyBalance, _ := balances.Load("ANY")
	expectedXESBalance := big.NewInt(0)
	expectedXESBalance.SetString("5483000000000000000000", 10)
	assert.Equal(t, expectedXESBalance, xesBalance)
	assert.Equal(t, big.NewInt(7373767001504), mkrBalance)
	assert.Equal(t, big.NewInt(0), anyBalance)
}
