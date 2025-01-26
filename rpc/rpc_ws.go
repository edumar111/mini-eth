// rpc/rpc_ws.go
package rpc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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

// Inicia el servidor WebSocket
func (srv *RPCServer) StartWS(port string) {
	http.HandleFunc("/ws", srv.wsHandler)
	log.Printf("Starting WS on port %s\n", port)
	// Nota: En producción, maneja mejor los errores de ListenAndServe
	http.ListenAndServe(":"+port, nil)
}

func (srv *RPCServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	go srv.handleConnection(conn)
}

func (srv *RPCServer) handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	for {
		// Leemos el mensaje entrante
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading WS message:", err)
			return
		}

		// Intentamos decodificar el mensaje a nuestra estructura RPCRequest
		var request RPCRequest
		if err := json.Unmarshal(msg, &request); err != nil {
			log.Println("Error unmarshaling request:", err)
			// Respondemos con un error de parseo
			srv.sendError(conn, -1, "Invalid JSON format")
			continue
		}

		// Creamos un objeto de respuesta
		response := RPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
		}

		// Revisamos el método
		switch request.Method {
		case "ping":
			// Si el método es "ping", respondemos con "pong"
			response.Result = "pong"

		default:
			// Si no se reconoce el método, enviamos un error
			response.Error = fmt.Sprintf("Method '%s' not found", request.Method)
		}

		// Enviamos la respuesta como JSON
		respBytes, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshaling response:", err)
			continue
		}

		// Mandamos el mensaje de vuelta al cliente WebSocket
		if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
			log.Println("Error writing WS message:", err)
			return
		}
	}
}

// Método auxiliar para enviar un error si el mensaje no es válido
func (srv *RPCServer) sendError(conn *websocket.Conn, id int, errMsg string) {
	response := RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   errMsg,
	}
	respBytes, _ := json.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, respBytes)
}
