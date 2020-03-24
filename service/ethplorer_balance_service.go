package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
)

type ethplorerBalanceService struct {
	smartContractTokensMap map[string]string
}

func NewEthplorerBalanceService(smartContractTokensMap map[string]string) *ethplorerBalanceService {
	return &ethplorerBalanceService{smartContractTokensMap: smartContractTokensMap}
}

func (me *ethplorerBalanceService) GetBalancesForAddress(ctx context.Context, address string) (*sync.Map, error) {
	resp, err := http.Get("http://api.ethplorer.io/getAddressInfo/" + address + "?apiKey=freekey")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	ethplorerResp := ethplorerResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ethplorerResp)
	if err != nil {
		return nil, err
	}

	balances := me.toMap(ethplorerResp)

	return balances, nil
}

func (me *ethplorerBalanceService) toMap(resp ethplorerResponse) *sync.Map {
	balances := new(sync.Map)
	ethFloat := big.NewFloat(resp.ETH.Balance)
	ethFloat.Mul(ethFloat, big.NewFloat(1000000000000000000))
	log.Println("Total?", ethFloat.String())
	ethString := ethFloat.Text('f', 0)
	ethBalance := big.NewInt(0)
	ethBalance.SetString(ethString, 10)
	balances.Store("ETH", ethBalance)

	tokensMap := me.responseTokensToMap(resp.Tokens)

	for _, tokenSymbol := range me.smartContractTokensMap {
		balance, found := tokensMap[tokenSymbol]
		if found {
			balances.Store(tokenSymbol, balance)
		} else {
			log.Printf("Token %s not found in %v", tokenSymbol, tokensMap)
			balances.Store(tokenSymbol, big.NewInt(0))
		}
	}

	return balances
}

func (me *ethplorerBalanceService) responseTokensToMap(tokens []token) map[string]*big.Int {
	balances := make(map[string]*big.Int)
	for _, token := range tokens {
		balances[token.Symbol] = token.Balance.BigIntValue
	}

	return balances
}

type BigInt struct {
	BigIntValue *big.Int
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	return []byte(b.BigIntValue.String()), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	stringValue := string(p)
	if stringValue == "null" {
		return nil
	}
	b.BigIntValue = new(big.Int)

	if strings.Contains(stringValue, "e") {
		// Scientific notation
		var newNum float64
		_, err := fmt.Sscanf(stringValue, "%e", &newNum)
		if err != nil {
			log.Printf("Can't convert %s, error: %v", stringValue, err)
			return err
		}

		b.BigIntValue.SetString(fmt.Sprintf("%.f", newNum), 10)
		return nil
	}

	_, ok := b.BigIntValue.SetString(stringValue, 10)
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", stringValue)
	}

	return nil
}

type ethplorerResponse struct {
	Address  string     `json:"address"`
	ETH      ethBalance `json:"ETH"`
	CountTxs int        `json:"countTxs"`
	Tokens   []token    `json:"tokens"`
}

type token struct {
	tokenInfo `json:"tokenInfo,omitempty"`
	Balance   BigInt `json:"balance"`
	TotalIn   int    `json:"totalIn"`
	TotalOut  int    `json:"totalOut"`
}

type ethBalance struct {
	Balance float64 `json:"balance"`
	Price   struct {
		Rate            float64 `json:"rate"`
		Diff            float64 `json:"diff"`
		Diff7D          float64 `json:"diff7d"`
		Ts              int     `json:"ts"`
		MarketCapUsd    float64 `json:"marketCapUsd"`
		AvailableSupply float64 `json:"availableSupply"`
		Volume24H       float64 `json:"volume24h"`
		Diff30D         float64 `json:"diff30d"`
	} `json:"price"`
}

type tokenInfo struct {
	Address           string `json:"address"`
	Name              string `json:"name"`
	Symbol            string `json:"symbol"`
	TotalSupply       string `json:"totalSupply"`
	Owner             string `json:"owner"`
	LastUpdated       int    `json:"lastUpdated"`
	IssuancesCount    int    `json:"issuancesCount"`
	HoldersCount      int    `json:"holdersCount"`
	EthTransfersCount int    `json:"ethTransfersCount"`
}
