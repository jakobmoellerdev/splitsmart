package splitcore_test

import (
	"slices"
	"testing"
	"time"

	"splitsmart/pkg/splitcore"
)

const dtFormat = "2006-01-02"

// testDate is a helper function to parse a date string into a time.Time
// "2020-01-01" -> time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).
func testDate(date string) time.Time {
	t, _ := time.Parse(dtFormat, date)

	return t
}

//nolint:funlen,maintidx
func TestGroupCalculations(t *testing.T) {
	cases := []struct {
		name                           string
		members                        splitcore.Members
		bills                          splitcore.Bills
		paymentsMade                   splitcore.Transactions
		expectedCalculatedTransactions splitcore.Transactions
	}{
		{
			name: "initial state",
			members: splitcore.Members{
				{ID: "Anne"},
				{ID: "Bob"},
				{ID: "Carla"},
			},
			bills: splitcore.Bills{
				{
					ID:                   "1",
					DistributionStrategy: splitcore.DistributionStrategyEqual,
					Date:                 testDate("2020-01-01"),
					PaidBy: map[splitcore.ID]splitcore.Amount{
						"Anne":  50, // 33.33 owed
						"Bob":   50, // 33.33 owed
						"Carla": 0,  // 33.33 owed
					},
				},
			},
			expectedCalculatedTransactions: splitcore.Transactions{
				{
					By:     "Carla",
					To:     "Anne",
					Amount: splitcore.Amount(float64(50+50+0) / float64(3) / float64(2)),
				},
				{
					By:     "Carla",
					To:     "Bob",
					Amount: splitcore.Amount(float64(50+50+0) / float64(3) / float64(2)),
				},
			},
		},
		{
			name: "initial state (with one payment made)",
			members: splitcore.Members{
				{ID: "Anne"},
				{ID: "Bob"},
				{ID: "Carla"},
			},
			bills: splitcore.Bills{
				{
					ID:                   "1",
					DistributionStrategy: splitcore.DistributionStrategyEqual,
					Date:                 testDate("2020-01-01"),
					PaidBy: map[splitcore.ID]splitcore.Amount{
						"Anne":  50, // 33.33 owed
						"Bob":   50, // 33.33 owed
						"Carla": 0,  // 33.33 owed
					},
				},
			},
			paymentsMade: splitcore.Transactions{
				{
					By:     "Carla",
					To:     "Anne",
					Amount: splitcore.Amount(float64(50+50+0) / float64(3) / float64(2)),
				},
			},
			expectedCalculatedTransactions: splitcore.Transactions{
				{
					By:     "Carla",
					To:     "Bob",
					Amount: splitcore.Amount(float64(50+50+0) / float64(3) / float64(2)),
				},
			},
		},
		{
			name: "simple 4 people split with one 1 person paying",
			members: splitcore.Members{
				{ID: "Anne"},
				{ID: "Bob"},
				{ID: "Carla"},
				{ID: "David"},
			},
			bills: splitcore.Bills{
				{
					ID:                   "1",
					DistributionStrategy: splitcore.DistributionStrategyEqual,
					Date:                 testDate("2020-01-01"),
					PaidBy: map[splitcore.ID]splitcore.Amount{
						"Anne":  100,
						"Bob":   0,
						"Carla": 0,
						"David": 0,
					},
				},
			},
			expectedCalculatedTransactions: splitcore.Transactions{
				{
					By:     "Bob",
					To:     "Anne",
					Amount: 25,
				},
				{
					By:     "Carla",
					To:     "Anne",
					Amount: 25,
				},
				{
					By:     "David",
					To:     "Anne",
					Amount: 25,
				},
			},
		},
		{
			name: "percentage based 4 people split with one 1 person paying",
			members: splitcore.Members{
				{ID: "Anne"},
				{ID: "Bob"},
				{ID: "Carla"},
				{ID: "David"},
			},
			bills: splitcore.Bills{
				{
					ID: "1",
					DistributionStrategy: splitcore.DistributionStrategyPercentage{
						"Anne":  0,
						"Bob":   0.25,
						"Carla": 0.25,
						"David": 0.5,
					},
					Date: testDate("2020-01-01"),
					PaidBy: map[splitcore.ID]splitcore.Amount{
						"Anne":  100,
						"Bob":   0,
						"Carla": 0,
						"David": 0,
					},
				},
			},
			expectedCalculatedTransactions: splitcore.Transactions{
				{
					By:     "Bob",
					To:     "Anne",
					Amount: 25,
				},
				{
					By:     "Carla",
					To:     "Anne",
					Amount: 25,
				},
				{
					By:     "David",
					To:     "Anne",
					Amount: 50,
				},
			},
		},
		{
			name: "complex 4 people splitting percentage scenario",
			members: splitcore.Members{
				{ID: "Anne"},
				{ID: "Bob"},
				{ID: "Carla"},
				{ID: "David"},
			},
			bills: splitcore.Bills{
				{
					ID: "1",
					DistributionStrategy: splitcore.DistributionStrategyPercentage{
						"Anne":  0.2,
						"Bob":   0.3,
						"Carla": 0.1,
						"David": 0.4,
					},
					Date: testDate("2020-01-01"),
					PaidBy: map[splitcore.ID]splitcore.Amount{
						"Anne":  100,
						"Bob":   0,
						"Carla": 0,
						"David": 0,
					},
				},
			},
			expectedCalculatedTransactions: splitcore.Transactions{
				{
					By:     "Bob",
					To:     "Anne",
					Amount: 30,
				},
				{
					By:     "Carla",
					To:     "Anne",
					Amount: 10,
				},
				{
					By:     "David",
					To:     "Anne",
					Amount: 40,
				},
			},
		},
		{
			name: "complex 4 people splitting shares scenario",
			members: splitcore.Members{
				{ID: "Anne"},
				{ID: "Bob"},
				{ID: "Carla"},
				{ID: "David"},
			},
			bills: splitcore.Bills{
				{
					ID: "1",
					DistributionStrategy: splitcore.DistributionStrategyShares{
						"Anne":  1,
						"Bob":   11,
						"Carla": 5,
						"David": 3,
					},
					Date: testDate("2020-01-01"),
					PaidBy: map[splitcore.ID]splitcore.Amount{
						"Anne":  100,
						"Bob":   0,
						"Carla": 0,
						"David": 0,
					},
				},
			},
			expectedCalculatedTransactions: splitcore.Transactions{
				{
					By:     "Bob",
					To:     "Anne",
					Amount: 55,
				},
				{
					By:     "Carla",
					To:     "Anne",
					Amount: 25,
				},
				{
					By:     "David",
					To:     "Anne",
					Amount: 15,
				},
			},
		},
		{
			name: "2 people paying scenario for 4 people splitting equally",
			members: splitcore.Members{
				{ID: "Anne"},
				{ID: "Bob"},
				{ID: "Carla"},
				{ID: "David"},
			},
			bills: splitcore.Bills{
				{
					ID:                   "1",
					DistributionStrategy: splitcore.DistributionStrategyEqual,
					Date:                 testDate("2020-01-01"),
					PaidBy: map[splitcore.ID]splitcore.Amount{
						"Anne":  50,
						"Bob":   50,
						"Carla": 0,
						"David": 0,
					},
				},
			},
			expectedCalculatedTransactions: splitcore.Transactions{
				{
					By:     "Carla",
					Amount: 25,
				},
				{
					By:     "David",
					Amount: 25,
				},
			},
		},
		{
			name: "2 people paying scenario for 4 people splitting percentage based",
			members: splitcore.Members{
				{ID: "Anne"},
				{ID: "Bob"},
				{ID: "Carla"},
				{ID: "David"},
			},
			bills: splitcore.Bills{
				{
					ID: "1",
					DistributionStrategy: splitcore.DistributionStrategyPercentage{
						"Anne":  0.1,
						"Bob":   0.2,
						"Carla": 0.2,
						"David": 0.5,
					},
					Date: testDate("2020-01-01"),
					PaidBy: map[splitcore.ID]splitcore.Amount{
						"Anne":  50,
						"Bob":   50,
						"Carla": 0,
						"David": 0,
					},
				},
			},
			expectedCalculatedTransactions: splitcore.Transactions{
				{
					By:     "David",
					To:     "Anne",
					Amount: 40,
				},
				{
					By:     "David",
					To:     "Bob",
					Amount: 10,
				},
				{
					By:     "Carla",
					To:     "Bob",
					Amount: 20,
				},
			},
		},
		{
			name: "4 people paying scenario for 4 people splitting percentage based",
			members: splitcore.Members{
				{ID: "Anne"},
				{ID: "Bob"},
				{ID: "Carla"},
				{ID: "David"},
			},
			bills: splitcore.Bills{
				{
					ID: "1",
					DistributionStrategy: splitcore.DistributionStrategyPercentage{
						"Anne":  0.1,
						"Bob":   0.2,
						"Carla": 0.2,
						"David": 0.5,
					},
					Date: testDate("2020-01-01"),
					PaidBy: map[splitcore.ID]splitcore.Amount{
						"Anne":  30,
						"Bob":   30,
						"Carla": 30,
						"David": 10,
					},
				},
			},
			expectedCalculatedTransactions: splitcore.Transactions{
				{
					By:     "David",
					To:     "Anne",
					Amount: 20,
				},
				{
					By:     "David",
					To:     "Bob",
					Amount: 10,
				},
				{
					By:     "David",
					To:     "Carla",
					Amount: 10,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			group := splitcore.NewGroup(tc.members)

			for _, bill := range tc.bills {
				group.AddBill(bill)
			}

			for _, payment := range tc.paymentsMade {
				group.AddPayment(payment)
			}

			transactions := group.CalculateTransactions()

			for _, expected := range tc.expectedCalculatedTransactions {
				if !slices.ContainsFunc(transactions, func(t splitcore.Transaction) bool {
					if t.By != expected.By {
						return false
					}
					if expected.To != "" && t.To != expected.To {
						return false
					}

					return t.Equals(expected.Amount)
				}) {
					t.Errorf("expected transaction not found: %v", expected)
				}
			}
		})
	}
}
