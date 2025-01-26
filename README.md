## mini-eth
Project to learn how to develop a blockchain node client
##
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

    curl -X POST --data '{"jsonrpc":"2.0","method":"ping","params":[],"id":1}' http://127.0.0.1:4045/rpc
WS
    wscat -c ws://localhost:4046/ws
    websocat  ws://localhost:4046/ws

{"jsonrpc":"2.0","method":"ping","params":[],"id":1}

{"jsonrpc":"2.0","error":"Method 'foo' not found","id":1}