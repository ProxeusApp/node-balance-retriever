module github.com/ProxeusApp/node-balance-retriever

go 1.13

require (
	github.com/ProxeusApp/proxeus-core v0.0.0-20200225075120-2ec19bc6c3a4
	github.com/ProxeusApp/proxeus-core/externalnode v0.0.0-20200224162123-278dfce8819f
	github.com/ethereum/go-ethereum v1.9.10
	github.com/labstack/echo v3.3.10+incompatible
	github.com/stretchr/testify v1.4.0
	github.com/valyala/fasttemplate v1.1.0 // indirect
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d // indirect
	golang.org/x/net v0.0.0-20200222125558-5a598a2470a0 // indirect
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
)

//replace github.com/ProxeusApp/proxeus-core/externalnode => ../proxeus-core/externalnode
