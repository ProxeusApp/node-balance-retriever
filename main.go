package main

import (
	"encoding/json"
	"errors"
	"github.com/ProxeusApp/proxeus-core/externalnode"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ProxeusApp/node-balance-retriever/service"
	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain/ethglue"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	serviceName = "balanceRetriever"
	jwtSecret   = "my secret 2"
	serviceUrl  = "127.0.0.1:8012"
	authKey     = "auth"
)

var (
	ethereumBalanceService service.EthereumBalanceService
	errCastingEthAddress   = errors.New("[taxreporter][next] casting error ethAddress")
)

func main() {
	ethClientUrl := os.Getenv("PROXEUS_ETH_CLIENT_URL")
	if len(ethClientUrl) == 0 {
		ethClientUrl = "https://ropsten.infura.io/v3/4876e0df8d31475799c8239ba2538c4c"
	}

	ethClient, err := ethglue.Dial(ethClientUrl)
	if err != nil {
		log.Println("[taxreporter][run] err: ", err.Error())
		return
	}

	xesAddress := os.Getenv("PROXEUS_XES_ADDRESS")
	if len(xesAddress) == 0 {
		//xesAddress = "0xA017ac5faC5941f95010b12570B812C974469c2C" //mainnet
		xesAddress = "0x84E0b37e8f5B4B86d5d299b0B0e33686405A3919" //ropsten
	}
	mkrAddress := os.Getenv("PROXEUS_MKR_ADDRESS")
	if len(mkrAddress) == 0 {
		//mkrAddress = "0x9f8F72aA9304c8B593d555F12eF6589cC3A579A2" //mainnet
		mkrAddress = "0x710129558e8fff5cab9c0c9c43b99d79ed864b99" //ropsten
	}

	tokensMap := map[string]string{
		xesAddress: "XES",
		mkrAddress: "MKR",
	}

	balanceService, err := service.NewEthClientBalanceService(ethClient, tokensMap)
	if err != nil {
		log.Println("[taxreporter][run] err: ", err.Error())
		return
	}

	ethereumBalanceService = service.NewEthereumBalanceService(balanceService)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.GET("/health", externalnode.Health)
	{
		g := e.Group("/node/:id")
		conf := middleware.DefaultJWTConfig
		conf.SigningKey = []byte(jwtSecret)
		conf.TokenLookup = "query:" + authKey
		g.Use(middleware.JWTWithConfig(conf))

		g.POST("/next", next)
		g.GET("/config", externalnode.Nop)
		g.POST("/config", externalnode.Nop)
		g.POST("/remove", externalnode.Nop)
		g.POST("/close", externalnode.Nop)
	}
	externalnode.Register(serviceName, serviceUrl, jwtSecret, "Converts currencies")
	err = e.Start(serviceUrl)
	if err != nil {
		log.Println("[taxreporter][run] err: ", err.Error())
	}
}

func next(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	ethAddress, ok := response["ethAddress"].(string)
	if !ok {
		return c.String(http.StatusInternalServerError, errCastingEthAddress.Error())
	}

	balanceResponse, err := ethereumBalanceService.GetBalances(c.Request().Context(), ethAddress)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	//fill response-map with balanceResponse result
	for k, v := range balanceResponse {
		response[k] = v.String()
	}

	return c.JSON(http.StatusOK, response)
}
