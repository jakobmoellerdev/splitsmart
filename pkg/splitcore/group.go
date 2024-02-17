package splitcore

import (
	"math"
	"slices"

	"github.com/google/uuid"
)

type (
	Group struct {
		Bills
		Members
		Transactions
		Balances
	}
	Balances map[ID]Amount
)

func NewGroup(members Members) *Group {
	return &Group{
		Members:      members,
		Bills:        make(Bills, 0),
		Transactions: make(Transactions, 0),
	}
}

func (grp *Group) AddPayment(transaction Transaction) {
	grp.Transactions = append(grp.Transactions, transaction)
	for _, member := range grp.Members {
		if member.ID == transaction.By {
			member.Sent = append(member.Sent, transaction)
		}

		if member.ID == transaction.To {
			member.Received = append(member.Received, transaction)
		}
	}

	grp.updateBalances()
}

func (grp *Group) RemovePayment(paymentID ID) {
	for i, payment := range grp.Transactions {
		if payment.ID == paymentID {
			grp.Transactions = slices.Concat(grp.Transactions[:i], grp.Transactions[i+1:])

			break
		}
	}

	for i, member := range grp.Members {
		for j, transaction := range member.Sent {
			if transaction.ID == paymentID {
				grp.Members[i].Sent = slices.Concat(member.Sent[:j], member.Sent[j+1:])

				break
			}
		}

		for j, transaction := range member.Received {
			if transaction.ID == paymentID {
				grp.Members[i].Received = slices.Concat(member.Received[:j], member.Received[j+1:])

				break
			}
		}
	}

	grp.updateBalances()
}

func (grp *Group) AddMember(member Member) {
	grp.Members = append(grp.Members, member)
}

func (grp *Group) RemoveMember(memberID ID) {
	for i, member := range grp.Members {
		if member.ID == memberID {
			grp.Members = slices.Concat(grp.Members[:i], grp.Members[i+1:])

			break
		}
	}

	for _, bill := range grp.Bills {
		if amount, ok := bill.PaidBy[memberID]; ok {
			delete(bill.PaidBy, memberID)
			// distribute the amount to the remaining members
			share := amount / Amount(len(bill.PaidBy))
			for member := range bill.PaidBy {
				bill.PaidBy[member] += share
			}
		}

		delete(bill.PaidBy, memberID)
	}

	grp.updateBalances()
}

func (grp *Group) AddBill(bill Bill) {
	// default to 0 for all members if not already involved
	for i := range grp.Members {
		if _, ok := bill.PaidBy[grp.Members[i].ID]; !ok {
			bill.PaidBy[grp.Members[i].ID] = 0
		}
	}

	grp.Bills = append(grp.Bills, bill)
	grp.updateBalances()
}

func (grp *Group) RemoveBill(billID ID) {
	for i, bill := range grp.Bills {
		if bill.ID == billID {
			grp.Bills = slices.Concat(grp.Bills[:i], grp.Bills[i+1:])

			break
		}
	}

	grp.updateBalances()
}

func (grp *Group) GetBalancesForBill(bill Bill) Balances {
	return bill.GetBalances(grp, bill)
}

func (grp *Group) updateBalances() {
	balances := make(Balances, len(grp.Members))

	for _, bill := range grp.Bills {
		balancesForBill := grp.GetBalancesForBill(bill)
		for member, amount := range balancesForBill {
			balances[member] += amount
		}
	}

	for _, settlement := range grp.Transactions {
		balances[settlement.By] -= settlement.Amount
		balances[settlement.To] += settlement.Amount
	}

	grp.Balances = balances
}

func (grp *Group) CalculateTransactions() Transactions {
	return grp.calculateTransactionsForBalances(grp.Balances)
}

//nolint:funlen
func (grp *Group) calculateTransactionsForBalances(balances Balances) Transactions {
	// Separate into debtors (to give money) and creditors (to be given money)
	debtors, creditors := Balances{}, Balances{}

	for member, amount := range balances {
		if amount < 0 {
			debtors[member] = amount
		} else if balances[member] > 0 {
			creditors[member] = amount
		}
	}

	type Balance struct {
		ID
		Amount
	}

	debtorsArray, creditorsArray := make([]Balance, 0), make([]Balance, 0)

	for member, amount := range debtors {
		debtorsArray = append(debtorsArray, Balance{member, amount})
	}

	for member, amount := range creditors {
		creditorsArray = append(creditorsArray, Balance{member, amount})
	}

	slices.SortFunc(debtorsArray, func(a, b Balance) int {
		return int(a.Amount - b.Amount)
	})

	slices.SortFunc(creditorsArray, func(a, b Balance) int {
		return int(b.Amount - a.Amount)
	})

	var transactions Transactions

	// Distribute debts from debtors to creditors
	for len(debtorsArray) > 0 && len(debtorsArray) > 0 {
		debtor := &debtorsArray[0]
		creditor := &creditorsArray[0]

		amount := Amount(math.Min(-float64(debtor.Amount), float64(creditor.Amount)))

		transactions = append(transactions, Transaction{
			ID:     ID(uuid.New().String()),
			By:     debtor.ID,
			To:     creditor.ID,
			Amount: amount,
		})

		// Update share
		debtor.Amount += amount
		creditor.Amount -= amount

		// Remove fully settled individuals
		if debtor.Amount.Equals(0) {
			debtorsArray = debtorsArray[1:]
		}

		if creditor.Amount.Equals(0) {
			creditorsArray = creditorsArray[1:]
		}
	}

	return transactions
}
