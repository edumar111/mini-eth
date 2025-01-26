// core/transaction.go
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
)

// Transaction representa una transacción simple.
type Transaction struct {
	From   string
	To     string
	Amount uint64
	Nonce  uint64
	// Para PoS podríamos agregar un campo `IsStake bool` o algo similar
	// Nonce, firma, etc., si quisiéramos mayor realismo
}

func (tx *Transaction) Hash() string {
	data := []byte(fmt.Sprintf("%s:%s:%d:%d", tx.From, tx.To, tx.Amount, tx.Nonce))
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

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
}
