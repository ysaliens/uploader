package handlers

import (
	"testing"
)

func TestCreateBudgetCodes(t *testing.T) {

	_, err := CreateBudgetCodes("wrong")
	if err == nil {
		t.Errorf("ERROR: Opened a missing file")
	}

	budgetCodeMap, err := CreateBudgetCodes("../files/config/budget_codes.csv")
	if err != nil || len(budgetCodeMap) == 0 {
		t.Errorf("ERROR: Budget Code hash map empty")
	}

}
