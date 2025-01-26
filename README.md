## mini-eth
Project to learn how to develop a blockchain node client
## Run


```shell script
./mini-eth run \
--p2p-port=30303 \
--rpc-http-port=4045 \
--rpc-ws-port=4046
```


```
mini-eth/
├── cmd/
│   └── mini-eth/        # Programa principal (CLI)
├── core/
│   ├── block.go         # Estructura y lógica de bloques
│   ├── transaction.go   # Estructura y lógica de transacciones
│   ├── genesis.go       # Estructura y lógica del bloque génesis
│   ├── state.go         # Manejo de estado, Merkle Trie, etc.
│   ├── consensus.go     # Lógica de 'stake' (o PoS muy simplificado)
│   └── token.go         # Lógica del token nativo
├── p2p/
│   ├── server.go        # Lógica del servidor P2P
│   └── peer.go          # Manejo de pares, conexión, mensajería
├── rpc/
│   ├── rpc_http.go      # Endpoints HTTP/JSON-RPC
│   └── rpc_ws.go        # Endpoints WS
├── cli/
│   └── commands.go      # Comandos de la CLI (start node, init genesis, etc.)
└── go.mod
```
RPC

    curl -X POST --data '{"jsonrpc":"2.0","method":"ping","params":[],"id":1}' http://127.0.0.1:4045
WS
    wscat -c ws://localhost:4046
    websocat  ws://localhost:4046

{"jsonrpc":"2.0","method":"ping","params":[],"id":1}

{"jsonrpc":"2.0","error":"Method 'foo' not found","id":1}

//
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0xc94770007dda54cF92009BFF0dE90c06F603a09f","latest"],"id":1}' http://127.0.0.1:4045