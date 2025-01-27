// core/state.go
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/cbergoon/merkletree" // ejemplo de librería Merkle Tree (3rd party)
	"sync"
)

// State representaría la estructura de balances y su Merkle Trie
type State struct {
	Balances   map[string]uint64
	Nonces     map[string]uint64
	merkleTree *merkletree.MerkleTree
	mu         sync.RWMutex
}

// Leaf implementa la interfaz merkletree.Content
type Leaf struct {
	Key   string
	Value uint64
}

func (l Leaf) CalculateHash() ([]byte, error) {
	h := sha256.Sum256([]byte(l.Key + string(rune(l.Value))))
	return h[:], nil
}
func (l Leaf) Equals(other merkletree.Content) (bool, error) {
	otherLeaf, ok := other.(Leaf)
	if !ok {
		return false, errors.New("type mismatch")
	}
	return l.Key == otherLeaf.Key && l.Value == otherLeaf.Value, nil
	//return l.Key == other.(Leaf).Key && l.Value == other.(Leaf).Value, nil
}

// NewState crea un state inicial
func NewState() *State {
	return &State{
		Balances:   make(map[string]uint64),
		Nonces:     make(map[string]uint64),
		merkleTree: nil,
	}
}

// Métodos para manipular Nonces
func (s *State) GetNonce(address string) uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Nonces[address]
}

func (s *State) IncrementNonce(address string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Nonces[address] = s.Nonces[address] + 1
}

// SetBalance establece un balance para una dirección
func (s *State) SetBalance(address string, amount uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Balances[address] = amount
}

// GetBalance obtiene el balance de una dirección
func (s *State) GetBalance(address string) uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Balances[address]
}

// AddTransaction ejemplo simplificado para añadir una transacción a "pendientes"
func (s *State) AddTransaction(tx *RawTx) {
	// Aquí podrías guardarlas en un pool de pendientes. Simplificado:
	s.IncrementNonce(tx.From()) // asumiendo que RawTx tenga un método From()
	balance := s.Balances[tx.From()]
	if balance >= tx.Value.Uint64() {
		s.Balances[tx.From()] -= tx.Value.Uint64()
		s.Balances[tx.To.Hex()] += tx.Value.Uint64()
	}
	s.UpdateMerkle()
}

// UpdateMerkle actualiza la Merkle Trie del State
func (s *State) UpdateMerkle() error {
	var list []merkletree.Content
	for k, v := range s.Balances {
		list = append(list, Leaf{Key: k, Value: v})
	}
	for k, v := range s.Nonces {
		list = append(list, Leaf{Key: "nonce_" + k, Value: v})
	}
	tree, err := merkletree.NewTree(list)
	if err != nil {
		return err
	}
	s.merkleTree = tree
	return nil
}

// Root devuelve la raíz Merkle del estado
func (s *State) Root() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.merkleTree == nil {
		return "", errors.New("Merkle Tree not initialized")
	}
	return hex.EncodeToString(s.merkleTree.MerkleRoot()), nil
}
