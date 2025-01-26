// core/token.go
package core

// Aquí definimos funciones de transferencia e inicialización de balances
func Transfer(balanceMap map[string]uint64, from, to string, amount uint64) bool {
	// Checar si `from` tiene saldo suficiente
	if balanceMap[from] < amount {
		return false
	}
	balanceMap[from] -= amount
	balanceMap[to] += amount
	return true
}
