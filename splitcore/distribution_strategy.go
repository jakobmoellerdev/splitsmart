package splitcore

type DistributionStrategy interface {
	GetBalances(group *Group, bill Bill) Balances
}

type (
	distributionStrategyEqual      struct{}
	DistributionStrategyPercentage map[ID]float64
	DistributionStrategyShares     map[ID]int
)

var DistributionStrategyEqual = distributionStrategyEqual{}

func (d distributionStrategyEqual) GetBalances(group *Group, bill Bill) Balances {
	// 50 / 2 = 25 for 2 people involved in the bill with 50 total
	share := bill.TotalAmount() / Amount(len(bill.PaidBy))
	distribution := make(Balances)

	for member, paid := range bill.PaidBy {
		distribution[member] = paid - share
	}

	return distribution
}

func (d DistributionStrategyPercentage) GetBalances(group *Group, bill Bill) Balances {
	distribution := make(Balances)
	total := bill.TotalAmount()

	for member, paid := range bill.PaidBy {
		distribution[member] = paid - (total * Amount(d[member]))
	}

	return distribution
}

func (d DistributionStrategyShares) GetBalances(group *Group, bill Bill) Balances {
	distribution := make(Balances)
	shareTotal := 0

	for _, share := range d {
		shareTotal += share
	}

	total := bill.TotalAmount()

	for member, paid := range bill.PaidBy {
		distribution[member] = paid - (total / Amount(shareTotal) * Amount(d[member]))
	}

	return distribution
}
