// core/transaction.go
package core

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"math/big"
	// go-ethereum libs (puedes reemplazarlas si prefieres otras)
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Transaction representa una transacción simple.
type RawTx struct {
	Nonce    uint64
	GasPrice *big.Int
	GasLimit *big.Int
	To       common.Address
	Value    *big.Int
	Data     []byte

	// Firma
	V *big.Int
	R *big.Int
	S *big.Int
}

// DecodedTx es lo que usaremos tras decodificar la TX RLP y verificar la firma
type DecodedTx struct {
	Nonce   uint64
	From    common.Address
	To      common.Address
	Value   *big.Int
	R, S, V *big.Int
}

// VerifySignature extrae la dirección `From` a partir de la firma (V,R,S)
func (tx *RawTx) VerifySignature() (common.Address, error) {
	// 1. Obtenemos el hash del "mensaje" a firmar (en Ethereum es la hash RLP + "Ethereum Signed Message" en caso de typed data, etc.)
	//    Para simplificar, supongamos que la firma es la transacción misma RLP-hasheada.
	//    En un escenario real, usarías algo como: sigHash = rlpHash(txWithoutSignature).
	sigHash := RawTxHash(tx)

	// 2. Reconstruimos la firma.
	//    Normalmente en Ethereum la "V" = 27/28 (o con chainId). Restamos 27 si es mayor que 28.
	v := tx.V.Uint64()
	if v == 27 || v == 28 {
		// OK
	} else if v >= 35 {
		// para EIP-155, etc., lo ignoramos en este ejemplo
		v -= 35
	}
	if v != 27 && v != 28 {
		return common.Address{}, errors.New("invalid signature (V)")
	}

	// Combine R, S, V en un signature de 65 bytes
	sig := make([]byte, 65)
	copy(sig[0:32], tx.R.Bytes())
	copy(sig[32:64], tx.S.Bytes())
	sig[64] = byte(v - 27) // 0 o 1

	// 3. Recuperamos la public key
	pubKey, err := crypto.Ecrecover(sigHash, sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pubKey) == 0 {
		return common.Address{}, errors.New("invalid pubkey from signature")
	}

	// 4. Convertimos pubkey en *ecdsa.PublicKey para obtener la address
	pubKeyECDSA, err := crypto.UnmarshalPubkey(pubKey)
	if err != nil {
		return common.Address{}, err
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKeyECDSA)
	return recoveredAddr, nil
}

// RawTxHash (simplificado) - en Ethereum se usa RLP. Aqui, una función mock
func RawTxHash(tx *RawTx) []byte {
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

// applyTxAndCreateBlock añade la TX a un nuevo bloque, lo aplica al State y actualiza el blockchain
func (tx *RawTx) ApplyTxAndCreateBlock(from common.Address, chain []*Block, state *State) (string, error) {
	// 1. Validar nonce
	currentNonce := state.GetNonce(from.Hex())
	fmt.Printf("From: %s \n", from.Hex())
	fmt.Printf("currentNonce %d \n", currentNonce)
	if tx.Nonce != currentNonce {
		return "", fmt.Errorf("invalid nonce: got %d, expected %d", tx.Nonce, currentNonce)
	}

	// 2. Validar balance
	fmt.Printf("From: %s Balance: %d \n", from.Hex(), state.Balances[from.Hex()])
	balance := state.Balances[from.Hex()]
	txValue := tx.Value.Uint64()
	fmt.Printf("txValue: %d \n", txValue)
	if balance < txValue {
		return "", fmt.Errorf("insufficient balance")
	}
	// 3. Aplicar transacción (transferencia)
	state.Balances[from.Hex()] = balance - txValue
	state.Balances[tx.To.Hex()] = state.Balances[tx.To.Hex()] + txValue
	state.IncrementNonce(from.Hex())
	state.UpdateMerkle()

	// 4. Crear un nuevo bloque con esta TX
	//newBlock := srv.createBlock([]core.DecodedTx{tx})

	stateRoot, _ := state.Root()

	newBlock := CreateBlock(chain, []RawTx{*tx}, stateRoot)
	// Lo añadimos a la "blockchain" (en tu caso, podrías tener un array de bloques)
	chain = append(chain, newBlock)

	// (Opcional) Llamar a tu mecanismo de consenso, broadcast, etc.
	log.Printf("New block created #%d with 1 TX\n", newBlock.Header.BlockNumber)

	txHash := common.BytesToHash(RawTxHash(tx))
	return txHash.Hex(), nil

}

/*
// Validar y aplicar transacciones
func (s *State) ApplyTransaction(tx *Transaction) error {
	// 1. Verificamos que el nonce sea el esperado
	currentNonce := s.GetNonce(tx.From)
	if tx.Nonce != currentNonce {
		return errors.New("invalid nonce")
	}

	// 2. Checamos que haya balance suficiente
	if s.Balances[tx.From] < tx.Amount {
		return errors.New("insufficient balance")
	}

	// 3. Efectuamos la transferencia
	s.Balances[tx.From] -= tx.Amount
	s.Balances[tx.To] += tx.Amount

	// 4. Incrementamos el nonce de la cuenta origen
	s.IncrementNonce(tx.From)

	// 5. (Opcional) Actualizar la Merkle Trie
	err := s.UpdateMerkle()
	if err != nil {
		return err
	}

	return nil
}*/
