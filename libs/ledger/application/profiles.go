package application

import (
	"fmt"
	"libs/ledger/application/commands"
	"libs/ledger/domain"
)

type DBSCreditCardProfile struct {
	ExpenseAccountID   int64
	LiabilityAccountID int64
}

func NewDBSImportProfile(expenseAccountID int64, liabilityAccountId int64) *DBSCreditCardProfile {
	return &DBSCreditCardProfile{
		ExpenseAccountID:   expenseAccountID,
		LiabilityAccountID: liabilityAccountId,
	}
}

func (p *DBSCreditCardProfile) Name() string { return "DBS Credit Card" }

func (p *DBSCreditCardProfile) NewPosting(raw commands.RawTransaction) (domain.Posting, error) {
	var posting domain.Posting
	var err error

	// NOTE: positive amount - money spent, negative amount - money refunded
	if raw.Amount >= 0 { // Increase Expense (Debit), Increase Liability (Credit)
		posting, err = domain.NewPosting(raw.ID, raw.Date, raw.Description, []domain.Entry{
			{AccountID: p.LiabilityAccountID, Description: raw.Description, DebitMicroSGD: 0, CreditMicroSGD: abs(raw.Amount)},
			{AccountID: p.ExpenseAccountID, Description: raw.Description, DebitMicroSGD: abs(raw.Amount), CreditMicroSGD: 0},
		})
	} else { // Refund: Decrease Liability (Debit) / Decrease Expense (Credit)
		posting, err = domain.NewPosting(raw.ID, raw.Date, raw.Description, []domain.Entry{
			{AccountID: p.LiabilityAccountID, Description: raw.Description, DebitMicroSGD: abs(raw.Amount), CreditMicroSGD: 0},
			{AccountID: p.ExpenseAccountID, Description: raw.Description, DebitMicroSGD: 0, CreditMicroSGD: abs(raw.Amount)},
		})
	}
	if err != nil {
		return domain.Posting{}, fmt.Errorf("error mapping raw transaction (%v): %w", raw, err)
	}

	return posting, nil
}

type OCBCStatementProfile struct {
	AssetsAccountID   int64
	ExpensesAccountID int64
	IncomeAccountID   int64
}

func NewOCBCStatementProfile(
	assetsAccountID int64,
	expensesAccountID int64,
	incomeAccountID int64,
) *OCBCStatementProfile {
	return &OCBCStatementProfile{
		AssetsAccountID:   assetsAccountID,
		ExpensesAccountID: expensesAccountID,
		IncomeAccountID:   incomeAccountID,
	}
}

func (p *OCBCStatementProfile) Name() string { return "DBS Credit Card" }

func (p *OCBCStatementProfile) NewPosting(raw commands.RawTransaction) (domain.Posting, error) {
	var posting domain.Posting
	var err error

	// NOTE: positive amount - income / salary, negative amount - money spent from bank account
	if raw.Amount >= 0 { // Increase Income (Debit), Increase Assets (Credit)
		posting, err = domain.NewPosting(raw.ID, raw.Date, raw.Description, []domain.Entry{
			{AccountID: p.IncomeAccountID, Description: raw.Description, DebitMicroSGD: 0, CreditMicroSGD: abs(raw.Amount)},
			{AccountID: p.AssetsAccountID, Description: raw.Description, DebitMicroSGD: abs(raw.Amount), CreditMicroSGD: 0},
		})
	} else { // Refund: Decrease Assets (Debit) / Increase Expense (Credit)
		posting, err = domain.NewPosting(raw.ID, raw.Date, raw.Description, []domain.Entry{
			{AccountID: p.ExpensesAccountID, Description: raw.Description, DebitMicroSGD: abs(raw.Amount), CreditMicroSGD: 0},
			{AccountID: p.AssetsAccountID, Description: raw.Description, DebitMicroSGD: 0, CreditMicroSGD: abs(raw.Amount)},
		})
	}
	if err != nil {
		return domain.Posting{}, fmt.Errorf("error mapping raw transaction (%v): %w", raw, err)
	}

	return posting, nil
}

func abs(i int64) int64 {
	if i < 0 {
		return -i
	}

	return i
}
