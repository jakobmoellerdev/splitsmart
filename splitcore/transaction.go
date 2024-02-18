package splitcore

import "fmt"

type (
	Transactions []Transaction
	Transaction  struct {
		ID
		By ID
		To ID
		Amount
	}
)

func (t Transaction) String() string {
	return fmt.Sprintf("Transaction(%s -> %v -> %s)", t.By, t.Amount, t.To)
}
