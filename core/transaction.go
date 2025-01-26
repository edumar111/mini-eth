// core/transaction.go
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Transaction representa una transacción simple.
type Transaction struct {
	From   string
	To     string
	Amount uint64
	// Para PoS podríamos agregar un campo `IsStake bool` o algo similar
	// Nonce, firma, etc., si quisiéramos mayor realismo
}

func (tx *Transaction) Hash() string {
	data := []byte(fmt.Sprintf("%s:%s:%d", tx.From, tx.To, tx.Amount))
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
