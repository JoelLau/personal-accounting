package packagesaccountingusecase

import (
	"testing"
)

func TestPackagesAccountingUseCase(t *testing.T) {
	result := PackagesAccountingUseCase("works")
	if result != "PackagesAccountingUseCase works" {
		t.Error("Expected PackagesAccountingUseCase to append 'works'")
	}
}
