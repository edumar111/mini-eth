// rpc/rpc_ws.go
package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/edumar111/my-geth-edu/core"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

// Estructuras para JSON-RPC
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	ID      int         `json:"id"`
}
type RPCWSServer struct {
	// Aquí también podemos almacenar un State o Blockchain
	State      *core.State
	Blockchain []*core.Block
}

// Inicia el servidor WebSocket
func (wsServer *RPCWSServer) StartWS(port string) {
	mux := http.NewServeMux()
	// Todas las peticiones a "/" en este puerto se tratarán como WebSocket
	mux.HandleFunc("/", wsServer.wsHandler)

	log.Printf("[WS-RPC] Listening on port %s (path /)\n", port)

	go func() {
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatalf("Error starting WS-RPC on port %s: %v", port, err)
		}
	}()
}
func (wsServer *RPCWSServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	go wsServer.handleConnection(conn)
}
func (wsServer *RPCWSServer) handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	for {
		// Leemos el mensaje entrante
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// Ver si es un error de cierre "normal" o uno "anormal".
			// isCloseError revisa el código de cierre.
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
				// Se pueden agregar otros que desees ignorar
			) {
				log.Printf("Conexión WS cerrada de forma normal: %v", err)
			} else {
				log.Printf("Error leyendo mensaje WS: %v", err)
			}
			return
		}

		var request RPCRequest
		if err := json.Unmarshal(msg, &request); err != nil {
			log.Println("Error unmarshaling request:", err)
			wsServer.sendError(conn, request.ID, "Invalid JSON format")
			continue
		}

		response := RPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
		}
		nodoRPC := &RPCServer{
			State:      wsServer.State,
			Blockchain: wsServer.Blockchain,
		}
		switch request.Method {
		case "ping":
			response.Result = "pong"

		case "eth_getTransactionCount":
			nonceHex, err := HandleGetTransactionCount(nodoRPC, request.Params)
			if err != nil {
				response.Error = err.Error()
			} else {
				response.Result = nonceHex
			}
		case "eth_sendRawTransaction":
			txHash, err := HandleSendRawTransaction(nodoRPC, request.Params)
			if err != nil {
				response.Error = err.Error()
			} else {
				response.Result = txHash
			}
		case "eth_getTransactionReceipt":
			txHash, err := HandleSendRawTransaction(nodoRPC, request.Params)
			if err != nil {
				response.Error = err.Error()
			} else {
				response.Result = txHash
			}
		default:
			response.Error = fmt.Sprintf("Method '%s' not found", request.Method)
		}

		respBytes, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshaling response:", err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
			log.Println("Error writing WS message:", err)
			return
		}
	}
}

// Método auxiliar para enviar un error si el mensaje no es válido
func (wsServer *RPCWSServer) sendError(conn *websocket.Conn, id int, errMsg string) {
	response := RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   errMsg,
	}
	respBytes, _ := json.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, respBytes)
}
