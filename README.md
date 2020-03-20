# Node Balance Retriever
An external node implementation for Proxeus core. Returns balances at a given date for ETH and erc20 tokens.

## Implementation

Current implementation uses a standard Ethereum node to calculate balances for Ether's balance + different ERC20 tokens.
Supported tokens: XES, MKR, BAT, OMG, ZRX, ENJ.

There's no caching and therefore **should only be used as demo purposes**.

Many requests to the Ethereum node will be made in order to calculate this data. 

## Usage

It is recommended to start it using docker.

The latest image is available at `proxeus/node-balance-retriever:latest`

See the configuration paragraph for more information on what environments variables can be overridden

## Configuration

The following parameters can be set via environment variables. 


| Environmentvariable | Required | Default value
--- | --- |   --- |  
PROXEUS_INFURA_API_KEY | X |  
PROXEUS_INSTANCE_URL |  | http://127.0.0.1:1323
SERVICE_NAME |  | Retrieve Token Balances
SERVICE_URL |  | http://localhost:SERVICE_PORT
SERVICE_PORT |  | 8012
SERVICE_SECRET |  | my secret 2
REGISTER_RETRY_INTERVAL |  | 5
PROXEUS_ETH_CLIENT_URL |  | https://ropsten.infura.io/v3/
PROXEUS_XES_ADDRESS |  | 0x84E0b37e8f5B4B86d5d299b0B0e33686405A3919
PROXEUS_MKR_ADDRESS |  | 0x710129558E8ffF5caB9c0c9c43b99d79Ed864B99
PROXEUS_BAT_ADDRESS |  | 0x60B10C134088ebD63f80766874e2Cade05fc987B
PROXEUS_OMG_ADDRESS |  | 0x9820B36a37Af9389a23ACfb7988C0ee6837763b6
PROXEUS_ZRX_ADDRESS |  | 0xA8E9Fa8f91e5Ae138C74648c9C304F1C75003A8D
PROXEUS_ENJ_ADDRESS |  | 0x81Ec0eD50441fc3d1d63763F27b24081E5b516d5

## Deployment

The node is available as docker image and can be used within a typical Proxeus Platform setup by including the following docker-compose service:

```
version: '3.7'

networks:
  xes-platform-network:
    name: xes-platform-network

services:
  node-balance-retriever:
    image: proxeus/node-balance-retriever:latest
    container_name: xes_node-node-balance-retriever
    networks:
      - xes-platform-network
    restart: unless-stopped
    environment:
      PROXEUS_INSTANCE_URL: http://xes-platform:1323
      PROXEUS_ETH_CLIENT_URL: "${PROXEUS_ETH_CLIENT_URL:-https://ropsten.infura.io/v3/}"
      PROXEUS_INFURA_API_KEY: ${PROXEUS_INFURA_API_KEY}
      SERVICE_SECRET: secret
      SERVICE_PORT: 8012
      SERVICE_URL: http://node-balance-retriever:8012
      TZ: Europe/Zurich
    ports:
      - "8012:8012"
```
