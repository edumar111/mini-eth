package rpc

import (
	"encoding/hex"
	"fmt"
	"github.com/edumar111/my-geth-edu/core"
	"github.com/ethereum/go-ethereum/rlp"
	"strconv"
	"strings"
)

// handleGetTransactionCount extrae el nonce y lo devuelve en hex
func HandleGetTransactionCount(srv *RPCServer, params []interface{}) (string, error) {
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

// handleSendRawTransaction decodifica la TX en hex RLP, la firma y la inserta en un nuevo bloque
func HandleSendRawTransaction(srv *RPCServer, params []interface{}) (string, error) {
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
	/*decodedTx := core.DecodedTx{
		Nonce: rawTx.Nonce,
		From:  fromAddr,
		To:    rawTx.To,
		Value: rawTx.Value,
		V:     rawTx.V,
		R:     rawTx.R,
		S:     rawTx.S,
	}*/

	// Aplicar la transacción al State
	// (aquí simplificamos: ignoramos gas, etc.)
	txHash, applyErr := rawTx.ApplyTxAndCreateBlock(fromAddr, srv.Blockchain, srv.State)
	if applyErr != nil {
		return "", applyErr
	}

	// Por convención, podríamos retornar un "txHash" que sería Keccak(rlp(tx))

	return txHash, nil
}
