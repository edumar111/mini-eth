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


curl -X POST -H "Content-Type: application/json" \
--data '{"jsonrpc":"2.0", "method":"eth_sendRawTransaction","params":["0xf869018203e882520894f17f52151ebef6c7334fad080c5704d77216b732881bc16d674ec80000801ba02da1c48b670996dcb1f447ef9ef00b33033c48a4fe938f420bec3e56bfd24071a062e0aa78a81bf0290afbc3a9d8e9a068e6d74caa66c5e0fa8a46deaae96b0833"],"id":1 }' \
http://127.0.0.1:4045/

curl -X POST -H "Content-Type: application/json" \
--data '{
"jsonrpc":"2.0",
"method":"eth_getTransactionReceipt",
"params":["0x6d6e4c616bd1bcaf2578e90fff8f3e7bd66eeb468196e9ee4df3745eb5c222e7"],
"id":53
}' \
http://127.0.0.1:4045