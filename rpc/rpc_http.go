// rpc/rpc_http.go
package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/edumar111/my-geth-edu/core"
	"log"
	"net/http"
	"strconv"
)

type RPCServer struct {
	// Podríamos tener referencia al objeto State, Blockchain, etc.
	// Referencia a State o a la blockchain. Aquí supongamos State directamente.
	State *core.State
}

// StartRPC arranca un servidor HTTP en el puerto indicado,
// y expone el handler en la ruta raíz "/".
func (srv *RPCServer) StartRPC(port string) {
	mux := http.NewServeMux()
	// Todas las peticiones que lleguen a "/" se procesan con handleRPC.
	mux.HandleFunc("/", srv.handleRPC)

	log.Printf("[HTTP-RPC] Listening on port %s (path /)\n", port)

	// Iniciamos el servidor con su propio mux
	go func() {
		err := http.ListenAndServe(":"+port, mux)
		if err != nil {
			log.Fatalf("Error starting HTTP-RPC on port %s: %v", port, err)
		}
	}()
}

// handleRPC ejemplo muy simple
// handleRPC procesa los requests JSON-RPC
func (srv *RPCServer) handleRPC(w http.ResponseWriter, r *http.Request) {
	var req RPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	response := RPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "ping":
		response.Result = "pong"
	case "eth_getTransactionCount":
		// Esperamos params[0] = address, params[1] = "latest"
		nonceHex, err := srv.handleGetTransactionCount(req.Params)
		if err != nil {
			response.Error = err.Error()
		} else {
			response.Result = nonceHex
		}

	// Otros métodos (ping, sendTransaction, etc.)

	default:
		response.Error = fmt.Sprintf("Method '%s' not found", req.Method)
	}

	respBytes, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

// handleGetTransactionCount extrae el nonce y lo devuelve en hex
func (srv *RPCServer) handleGetTransactionCount(params []interface{}) (string, error) {
	// Validamos parámetros
	if len(params) < 2 {
		return "", fmt.Errorf("invalid params")
	}
	address, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid address param")
	}
	blockParam, ok := params[1].(string)
	if !ok {
		return "", fmt.Errorf("invalid block param")
	}
	if blockParam != "latest" {
		return "", fmt.Errorf("only 'latest' blockParam is supported")
	}

	nonce := srv.State.GetNonce(address)
	// Lo retornamos en formato hex '0x...' como hace Ethereum
	nonceHex := "0x" + strconv.FormatUint(nonce, 16)
	return nonceHex, nil
}
