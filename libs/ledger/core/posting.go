package core

import (
	"errors"
	"fmt"
	"time"
)

type Posting struct {
	// WARN: this is **NOT** the database ID
	// unique identifier to check if posting has already been processed
	ID          string
	Date        time.Time // when the transaction was made
	Description string    // name of the posting
	Entries     []Entry   // list of how we will track this posting
}

func NewPosting(id string, date time.Time, description string,
	entries []Entry) (Posting, error) {

	p := Posting{
		ID:          id,
		Date:        date,
		Description: description,
		Entries:     entries,
	}

	return p, p.Error()
}

func (p *Posting) DebitMicroSGD() int64 {
	var debit int64 = 0
	for _, e := range p.Entries {
		debit += e.DebitMicroSGD
	}
	return debit
}

func (p *Posting) CreditMicroSGD() int64 {
	var credit int64 = 0
	for _, e := range p.Entries {
		credit += e.DebitMicroSGD
	}
	return credit
}

func (p *Posting) Error() error {
	var errs = []error{}

	for _, e := range p.Entries {
		if e.DebitMicroSGD < 0 || e.CreditMicroSGD < 0 {
			errs = append(errs, fmt.Errorf("amounts cannot be negative: %v", e))
		}
		if e.DebitMicroSGD > 0 && e.CreditMicroSGD > 0 {
			errs = append(errs, fmt.Errorf("entry cannot be both debit and credit: %v", e))
		}
	}

	debit := p.DebitMicroSGD()
	credit := p.CreditMicroSGD()
	if debit != credit {
		errs = append(errs, fmt.Errorf("unbalanced: debit %d != credit %d", debit, credit))
	}

	return errors.Join(errs...)
}
