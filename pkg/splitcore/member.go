package splitcore

type (
	Members []Member
	Member  struct {
		ID
		History
	}

	History struct {
		Received Transactions
		Sent     Transactions
	}
)
