package main

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"fmt"
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
	serviceName       = "balanceRetriever"
	defaultServiceUrl = "127.0.0.1:8012"
	jwtSecret         = "my secret 2"
	defaultProxeusUrl = "http://127.0.0.1:1323"
	defaultAuthkey    = "auth"
)

var (
	ethereumBalanceService service.EthereumBalanceService
	errCastingEthAddress   = errors.New("[taxreporter][next] casting error ethAddress")
)

func main() {
	proxeusUrl := os.Getenv("PROXEUS_INSTANCE_URL")
	if len(proxeusUrl) == 0 {
		proxeusUrl = defaultProxeusUrl
	}
	serviceUrl := os.Getenv("SERVICE_URL")
	if len(serviceUrl) == 0 {
		serviceUrl = defaultServiceUrl
	}
	authKey := os.Getenv("AUTH_KEY")
	if len(authKey) == 0 {
		authKey = defaultAuthkey
	}
	fmt.Println()
	fmt.Println("#######################################################")
	fmt.Println("# STARTING NODE - " + serviceName)
	fmt.Println("# listing on " + serviceUrl)
	fmt.Println("# connecting to " + proxeusUrl)
	fmt.Println("#######################################################")
	fmt.Println()
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
		mkrAddress = "0x710129558E8ffF5caB9c0c9c43b99d79Ed864B99" //ropsten
	}

	batAddress := os.Getenv("PROXEUS_BAT_ADDRESS")
	if len(batAddress) == 0 {
		//batAddress = "" //mainnet
		batAddress = "0x60B10C134088ebD63f80766874e2Cade05fc987B" //ropsten
	}
	usdcAddress := os.Getenv("PROXEUS_USDC_ADDRESS")
	if len(usdcAddress) == 0 {
		//usdcAddress = "" //mainnet
		usdcAddress = "0xFE724a829fdF12F7012365dB98730EEe33742ea2" //ropsten
	}
	repAddress := os.Getenv("PROXEUS_REP_ADDRESS")
	if len(repAddress) == 0 {
		//repAddress = "" //mainnet
		repAddress = "0xc853bA17650D32DAba343294998eA4E33e7a48B9" //ropsten
	}
	omgAddress := os.Getenv("PROXEUS_OMG_ADDRESS")
	if len(omgAddress) == 0 {
		//omgAddress = "" //mainnet
		omgAddress = "0x9820B36a37Af9389a23ACfb7988C0ee6837763b6" //ropsten
	}
	linkAddress := os.Getenv("PROXEUS_LINK_ADDRESS")
	if len(linkAddress) == 0 {
		//linkAddress = "" //mainnet
		linkAddress = "0x20fE562d797A42Dcb3399062AE9546cd06f63280" //ropsten
	}
	zrxAddress := os.Getenv("PROXEUS_ZRX_ADDRESS")
	if len(zrxAddress) == 0 {
		//zrxAddress = "" //mainnet
		zrxAddress = "0xA8E9Fa8f91e5Ae138C74648c9C304F1C75003A8D" //ropsten
	}
	enjAddress := os.Getenv("PROXEUS_ENJ_ADDRESS")
	if len(enjAddress) == 0 {
		//enjAddress = "" //mainnet
		enjAddress = "0x81Ec0eD50441fc3d1d63763F27b24081E5b516d5" //ropsten
	}

	//make sure to add new contract addresses with checksum (EIP-55)
	tokensMap := map[string]string{
		common.HexToAddress(xesAddress).String():  "XES",
		common.HexToAddress(mkrAddress).String():  "MKR",
		common.HexToAddress(batAddress).String():  "BAT",
		common.HexToAddress(usdcAddress).String(): "USDC",
		common.HexToAddress(repAddress).String():  "REP",
		common.HexToAddress(omgAddress).String():  "OMG",
		common.HexToAddress(linkAddress).String(): "LINK",
		common.HexToAddress(zrxAddress).String():  "ZRX",
		common.HexToAddress(enjAddress).String():  "ENJ",
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
	externalnode.Register(proxeusUrl, serviceName, serviceUrl, jwtSecret, "Retrieves token balances of an address")
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
