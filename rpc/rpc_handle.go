package rpc

import (
	"encoding/hex"
	"fmt"
	"github.com/edumar111/my-geth-edu/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
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

// HandleSendRawTransaction decodifica la TX en hex RLP, la firma y la inserta en un nuevo bloque
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

	// Aplicar la transacción al State
	// (aquí simplificamos: ignoramos gas, etc.)
	txHash, applyErr := rawTx.ApplyTxAndCreateBlock(fromAddr, &srv.Blockchain, srv.State)
	fmt.Printf("blockchain final %v+", srv.Blockchain)
	if applyErr != nil {
		return "", applyErr
	}

	// Por convención, podríamos retornar un "txHash" que sería Keccak(rlp(tx))

	return txHash, nil
}

// HandleGetTransactionReceipt busca la TX por hash y retorna un objeto JSON
// con blockNumber, transactionIndex y los campos de RawTx.
func HandleGetTransactionReceipt(srv *RPCServer, params []interface{}) (interface{}, error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("missing tx hash param")
	}
	hashParam, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid param type for tx hash")
	}
	// Remover el "0x" si existe
	hashParam = strings.TrimPrefix(hashParam, "0x")

	// Convertir el hash a bytes
	txHashBytes, err := hex.DecodeString(hashParam)
	if err != nil {
		return nil, fmt.Errorf("invalid hex tx hash: %v", err)
	}

	// Buscamos la TX en la blockchain
	blockIndex, txIndex, rawTx := findTransactionByHash(srv, txHashBytes)
	if rawTx == nil {
		return nil, fmt.Errorf("transaction not found")
	}

	// Obtenemos blockNumber en hex
	blockNumHex := "0x" + strconv.FormatUint(srv.Blockchain[blockIndex].Header.BlockNumber, 16)
	// Obtenemos transactionIndex en hex
	txIndexHex := "0x" + strconv.FormatUint(uint64(txIndex), 16)

	// Convertir valores a hex
	nonceHex := "0x" + strconv.FormatUint(rawTx.Nonce, 16)
	gasPriceHex := bigIntToHex(rawTx.GasPrice)
	gasLimitHex := bigIntToHex(rawTx.GasLimit)
	valueHex := bigIntToHex(rawTx.Value)
	vHex := bigIntToHex(rawTx.V)
	rHex := bigIntToHex(rawTx.R)
	sHex := bigIntToHex(rawTx.S)

	// Address en hex (ej. "0xabc123..."), Data en hex
	toHex := rawTx.To.Hex()
	dataHex := "0x" + hex.EncodeToString(rawTx.Data)

	// Construimos el objeto que retornaremos
	receipt := map[string]interface{}{
		"blockNumber":      blockNumHex,
		"transactionIndex": txIndexHex,

		"nonce":    nonceHex,
		"gasPrice": gasPriceHex,
		"gasLimit": gasLimitHex,
		"to":       toHex,
		"value":    valueHex,
		"data":     dataHex,

		"v": vHex,
		"r": rHex,
		"s": sHex,
	}

	// Si en algún punto calculas "from" mediante firma, puedes agregarlo:
	// "from": "0x..."
	// (Por defecto Ethereum no lo decodifica del RLP, sino de la firma.)

	return receipt, nil
}
func bigIntToHex(x *big.Int) string {
	if x == nil {
		return "0x0"
	}
	// Si x = 0 => "0x0"
	if x.Sign() == 0 {
		return "0x0"
	}
	// Si x > 0 => "0x<hex>"
	return "0x" + x.Text(16)
}
func findTransactionByHash(srv *RPCServer, hashBytes []byte) (int, int, *core.RawTx) {
	wantedHash := common.BytesToHash(hashBytes)

	fmt.Println("Find wantedHash", wantedHash)
	fmt.Printf("Blockchain size: %d \n", len(srv.Blockchain))
	for bIndex, block := range srv.Blockchain {
		fmt.Println("Find block", bIndex)
		fmt.Printf("Transactions size: %d \n", len(block.Transactions))
		for tIndex, tx := range block.Transactions {
			txHash := common.BytesToHash(RawTxHash(tx))
			//txHash := tx.
			fmt.Println("found txHash ", txHash)
			if txHash == wantedHash {
				return bIndex, tIndex, tx
			}
		}
	}
	return -1, -1, nil
}
func RawTxHash(tx *core.RawTx) []byte {
	// O usar rlp.EncodeToBytes(...)
	// Para ejemplo, supongamos un keccak256 de:
	// nonce + to + value
	toBytes := tx.To.Bytes()
	data := append(
		append(
			common.BigToHash(new(big.Int).SetUint64(tx.Nonce)).Bytes(),
			toBytes...,
		),
		tx.Value.Bytes()...,
	)
	h := crypto.Keccak256Hash(data)
	return h.Bytes()
}

/*
func calculateTxHash(rawTx *core.RawTx) common.Hash {

	encoded, _ := rlp.EncodeToBytes(rawTx)
	return common.BytesToHash(crypto.Keccak256(encoded))
}*/
