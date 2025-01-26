package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Estructuras para la petición y respuesta JSON-RPC
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

func main() {
	// 1. Generar una clave privada y dirección
	privKey, address, err := generateWallet()
	if err != nil {
		log.Fatalf("Error generating wallet: %v", err)
	}

	fmt.Println("Private Key (hex):", hex.EncodeToString(privKey.D.Bytes()))
	fmt.Println("Address:", address)

	// 2. Invocar al método eth_getTransactionCount en nuestro nodo
	rpcURL := "http://localhost:4045" // Ajusta la URL/puerto según tu configuración
	nonceHex, err := getTransactionCount(rpcURL, address)
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	fmt.Printf("Nonce de la cuenta %s: %s\n", address, nonceHex)
}

// generateWallet genera un par de llaves ECDSA (secp256k1) y produce una dirección Ethereum.
func generateWallet() (*ecdsa.PrivateKey, string, error) {
	// 1. Genera la clave privada
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// Nota: Ethereum usa secp256k1, que en Go puedes usar a través de la librería "github.com/ethereum/go-ethereum/crypto",
	// pero en este ejemplo didáctico usamos elliptic.P256() del estándar.
	if err != nil {
		return nil, "", err
	}

	// 2. Obtenemos la clave pública en formato (x, y)
	pubKey := privKey.PublicKey

	// 3. Serializamos la parte X,Y de la pubKey para luego hacerle keccak256
	//    (En Ethereum es la parte sin 0x04 y tomamos los últimos 20 bytes del hash)
	//    Aquí hacemos un ejemplo muy simplificado usando P256. Con secp256k1 es similar.
	pubBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)

	// 4. Usamos keccak256 para obtener la dirección
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes[1:]) // A veces se omite el primer byte 0x04 en secp256k1
	hashed := hash.Sum(nil)

	// Tomamos los últimos 20 bytes
	addressBytes := hashed[len(hashed)-20:]
	address := "0x" + hex.EncodeToString(addressBytes)

	// Lo pasamos a minúsculas para mantener estilo Ethereum
	address = strings.ToLower(address)

	return privKey, address, nil
}

// getTransactionCount realiza un POST JSON-RPC a eth_getTransactionCount
func getTransactionCount(rpcURL string, address string) (string, error) {
	// Preparamos la petición JSON-RPC
	reqBody := RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getTransactionCount",
		Params:  []interface{}{address, "latest"},
		ID:      1,
	}

	// Convertimos a JSON
	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Hacemos el request HTTP POST
	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parseamos la respuesta
	var rpcResp RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return "", err
	}

	// Si hay error en la respuesta JSON-RPC, lo mostramos
	if rpcResp.Error != nil {
		return "", fmt.Errorf("RPC Error: %v", rpcResp.Error)
	}

	// rpcResp.Result debería contener el nonce en formato "0x..." (hex)
	nonce, ok := rpcResp.Result.(string)
	if !ok {
		return "", fmt.Errorf("cannot parse nonce result")
	}

	return nonce, nil
}
