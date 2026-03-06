package core_test

import (
	core "packages/accounting/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPosting(t *testing.T) {
	t.Parallel()

	var (
		accountID__Liabilities_CreditCard int64 = 2001
		accountID__Expenses_EatingOut     int64 = 4402
		accountID__Expenses_Gifts         int64 = 4302
	)

	var (
		totalBill   int64 = 100_000_000           // 100.00 SGD
		myPortion   int64 = totalBill / 4         //  25.00 SGD
		giftPortion int64 = totalBill - myPortion //  75.00 SGD
	)

	posting := core.Posting{
		ID:          "4119-123412-123|06/03/2026|06/03/2026|McDonald's|SGD 100.00|", // example csv line
		Description: "McDonald's",
		Date:        time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC),
		Entries: []core.Entry{
			{AccountID: accountID__Liabilities_CreditCard, Description: "McDonald's", DebitMicroSGD: 0, CreditMicroSGD: totalBill},
			{AccountID: accountID__Expenses_EatingOut, Description: "McDonald's", DebitMicroSGD: myPortion, CreditMicroSGD: 0},
			{AccountID: accountID__Expenses_Gifts, Description: "McDonald's", DebitMicroSGD: giftPortion, CreditMicroSGD: 0},
		},
	}

	// balanced, as all things should be
	assert.Equal(t, posting.DebitMicroSGD(), totalBill)
	assert.Equal(t, posting.CreditMicroSGD(), totalBill)
}
