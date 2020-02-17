package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ProxeusApp/proxeus-core/externalnode"

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

	batAddress := os.Getenv("PROXEUS_BAT_ADDRESS")
	if len(batAddress) == 0 {
		//batAddress = "" //mainnet
		batAddress = "0x60b10c134088ebd63f80766874e2cade05fc987b" //ropsten
	}
	usdcAddress := os.Getenv("PROXEUS_USDC_ADDRESS")
	if len(usdcAddress) == 0 {
		//usdcAddress = "" //mainnet
		usdcAddress = "0xfe724a829fdf12f7012365db98730eee33742ea2" //ropsten
	}
	repAddress := os.Getenv("PROXEUS_REP_ADDRESS")
	if len(repAddress) == 0 {
		//repAddress = "" //mainnet
		repAddress = "0xc853ba17650d32daba343294998ea4e33e7a48b9" //ropsten
	}
	omgAddress := os.Getenv("PROXEUS_OMG_ADDRESS")
	if len(omgAddress) == 0 {
		//omgAddress = "" //mainnet
		omgAddress = "0x9820b36a37af9389a23acfb7988c0ee6837763b6" //ropsten
	}
	linkAddress := os.Getenv("PROXEUS_LINK_ADDRESS")
	if len(linkAddress) == 0 {
		//linkAddress = "" //mainnet
		linkAddress = "0x20fe562d797a42dcb3399062ae9546cd06f63280" //ropsten
	}
	zrxAddress := os.Getenv("PROXEUS_ZRX_ADDRESS")
	if len(zrxAddress) == 0 {
		//zrxAddress = "" //mainnet
		zrxAddress = "0xa8e9fa8f91e5ae138c74648c9c304f1c75003a8d" //ropsten
	}
	enjAddress := os.Getenv("PROXEUS_ENJ_ADDRESS")
	if len(enjAddress) == 0 {
		//enjAddress = "" //mainnet
		enjAddress = "0x81ec0ed50441fc3d1d63763f27b24081e5b516d5" //ropsten
	}

	tokensMap := map[string]string{
		xesAddress:  "XES",
		mkrAddress:  "MKR",
		batAddress:  "BAT",
		usdcAddress: "USDC",
		repAddress:  "REP",
		omgAddress:  "OMG",
		linkAddress: "LINK",
		zrxAddress:  "ZRX",
		enjAddress:  "ENJ",
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
