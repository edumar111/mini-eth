// rpc/rpc_http.go
package rpc

import (
	"encoding/json"
	"log"
	"net/http"
)

type RPCServer struct {
	// Podríamos tener referencia al objeto State, Blockchain, etc.
}

func (srv *RPCServer) StartRPC(port string) {
	http.HandleFunc("/rpc", srv.handleRPC)
	log.Printf("Starting RPC on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

// handleRPC ejemplo muy simple
func (srv *RPCServer) handleRPC(w http.ResponseWriter, r *http.Request) {
	type RPCRequest struct {
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
	}
	var req RPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Demo: respondemos con un "pong" si method == "ping"
	switch req.Method {
	case "ping":
		json.NewEncoder(w).Encode(map[string]string{"result": "pong"})
	// Implementar otros métodos (ej. "sendTransaction", "getBalance", etc.)
	default:
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not found"})
	}
}
