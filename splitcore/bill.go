package splitcore

import "time"

type (
	Bills []Bill
	Bill  struct {
		ID
		DistributionStrategy
		Date   time.Time
		PaidBy map[ID]Amount
	}
)

func (bill Bill) TotalAmount() Amount {
	total := Amount(0.0)
	for _, amount := range bill.PaidBy {
		total += amount
	}

	return total
}
