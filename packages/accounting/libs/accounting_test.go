package accounting_test

import (
	"testing"
)

func TestLibsAccounting(t *testing.T) {
	result := accounting.LibsAccounting("works")
	if result != "LibsAccounting works" {
		t.Error("Expected LibsAccounting to append 'works'")
	}
}
