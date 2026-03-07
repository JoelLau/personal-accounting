package domain

// NOTE: 1 SGD = 1,000,000 mSGD
type Entry struct {
	// reference to ledger account ID
	AccountID      int64
	Description    string
	DebitMicroSGD  int64
	CreditMicroSGD int64
}
