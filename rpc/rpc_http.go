// rpc/rpc_http.go
package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/edumar111/my-geth-edu/core"
	"log"
	"net/http"
)

type RPCServer struct {
	// Podríamos tener referencia al objeto State, Blockchain, etc.
	// Referencia a State o a la blockchain. Aquí supongamos State directamente.
	State      *core.State
	Blockchain []*core.Block
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
		nonceHex, err := HandleGetTransactionCount(srv, req.Params)
		if err != nil {
			response.Error = err.Error()
		} else {
			response.Result = nonceHex
		}
	case "eth_sendRawTransaction":
		txHash, err := HandleSendRawTransaction(srv, req.Params)
		if err != nil {
			response.Error = err.Error()
		} else {
			response.Result = txHash // Retornamos un "txHash" por ejemplo
		}
	case "eth_getTransactionReceipt":
		txHash, err := HandleGetTransactionReceipt(srv, req.Params)
		if err != nil {
			response.Error = err.Error()
		} else {
			response.Result = txHash // Retornamos un "txHash" por ejemplo
		}
	// Otros métodos (ping, sendTransaction, etc.)

	default:
		response.Error = fmt.Sprintf("Method '%s' not found", req.Method)
	}

	respBytes, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

// handleSendRawTransaction decodifica la TX en hex RLP, la firma y la inserta en un nuevo bloque
/*func (srv *RPCServer) handleSendRawTransaction(params []interface{}) (string, error) {
	//TODO solo para probar la demo****************
	srv.State.Nonces["0x4709421B04e70e3925dFC86307727b588709C7bB"] = 1
	srv.State.Balances["0x4709421B04e70e3925dFC86307727b588709C7bB"] = 9000000000000000000 //9 ETH
	// Esperamos un array con 1 string en hex
	if len(params) < 1 {
		return "", fmt.Errorf("missing parameter")
	}
	rawHex, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid param type")
	}

	// Removemos "0x" si existe
	rawHex = strings.TrimPrefix(rawHex, "0x")

	// Decodificamos hex a bytes
	rawBytes, err := hex.DecodeString(rawHex)
	if err != nil {
		return "", fmt.Errorf("hex decode error: %v", err)
	}

	// Decodificar via RLP
	var rawTx core.RawTx
	err = rlp.DecodeBytes(rawBytes, &rawTx)
	if err != nil {
		return "", fmt.Errorf("RLP decode error: %v", err)
	}

	// Verificar firma y extraer "from" address
	fromAddr, err := rawTx.VerifySignature()
	if err != nil {
		return "", fmt.Errorf("signature verify error: %v", err)
	}

	// Crear una transacción "interna" en nuestro State


	// Aplicar la transacción al State
	// (aquí simplificamos: ignoramos gas, etc.)
	txHash, applyErr := rawTx.ApplyTxAndCreateBlock(fromAddr, srv.Blockchain, srv.State)
	if applyErr != nil {
		return "", applyErr
	}

	// Por convención, podríamos retornar un "txHash" que sería Keccak(rlp(tx))

	return txHash, nil
}*/

// / que nos handles
// applyTxAndCreateBlock añade la TX a un nuevo bloque, lo aplica al State y actualiza el blockchain
/*func (srv *RPCServer) applyTxAndCreateBlock(tx core.DecodedTx) error {
	// 1. Validar nonce
	currentNonce := srv.State.GetNonce(tx.From.Hex())
	fmt.Printf("From: %s \n", tx.From.Hex())
	fmt.Printf("currentNonce %d \n", currentNonce)
	if tx.Nonce != currentNonce {
		return fmt.Errorf("invalid nonce: got %d, expected %d", tx.Nonce, currentNonce)
	}

	// 2. Validar balance
	fmt.Printf("From: %s Balance: %d \n", tx.From.Hex(), srv.State.Balances[tx.From.Hex()])
	balance := srv.State.Balances[tx.From.Hex()]
	txValue := tx.Value.Uint64()
	fmt.Printf("txValue: %d \n", txValue)
	if balance < txValue {
		return fmt.Errorf("insufficient balance")
	}
	// 3. Aplicar transacción (transferencia)
	srv.State.Balances[tx.From.Hex()] = balance - txValue
	srv.State.Balances[tx.To.Hex()] = srv.State.Balances[tx.To.Hex()] + txValue
	srv.State.IncrementNonce(tx.From.Hex())
	srv.State.UpdateMerkle()

	// 4. Crear un nuevo bloque con esta TX
	//newBlock := srv.createBlock([]core.DecodedTx{tx})
	rawTx := core.RawTx{
		To:    tx.To,
		Value: tx.Value,
		Nonce: tx.Nonce,
		V:     tx.V,
		R:     tx.R,
		S:     tx.S,
	}
	stateRoot, _ := srv.State.Root()
	newBlock := core.CreateBlock(srv.Blockchain, []core.RawTx{rawTx}, stateRoot)
	// Lo añadimos a la "blockchain" (en tu caso, podrías tener un array de bloques)
	srv.Blockchain = append(srv.Blockchain, newBlock)

	// (Opcional) Llamar a tu mecanismo de consenso, broadcast, etc.
	log.Printf("New block created #%d with 1 TX\n", newBlock.Header.BlockNumber)
	return nil
}*/

// createBlock ejemplo simplificado
/*
func (srv *RPCServer) createBlock(rawTxs []core.RawTx) *core.Block {

	var txPointers []*core.RawTx
	for _, rt := range rawTxs {
		txCopy := rt // para evitar issues de range
		txPointers = append(txPointers, &txCopy)
	}

	parentHash := "0x0"
	if len(srv.Blockchain) > 0 {
		parentHash = srv.Blockchain[len(srv.Blockchain)-1].Hash()
	}
	blockNumber := uint64(len(srv.Blockchain))

	// Obtener stateRoot si lo necesitas
	stateRoot, _ := srv.State.Root()

	// Crear el bloque
	block := core.NewBlock(parentHash, blockNumber, txPointers, stateRoot)
	// Anexar el bloque a tu chain
	srv.Blockchain = append(srv.Blockchain, block)
	return block
}*/
