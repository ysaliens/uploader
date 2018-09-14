package handlers

import (
	"encoding/csv"
	"io"
	"os"
)

// BudgetCode struct for hash map
type BudgetCode struct {
	Opex     string
	Category string
	Label    string
}

// CreateBudgetCodes - Read budget codes CSV and create a hash map
func CreateBudgetCodes(file string) (map[string]*BudgetCode, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	budgetCodes := map[string]*BudgetCode{}
	csvr := csv.NewReader(f)
	for { // Streaming CSV reader to minimize memory
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return budgetCodes, err
		}

		bc := &BudgetCode{}
		bc.Opex = row[0]
		bc.Category = row[1]
		bc.Label = row[3]
		budgetCodes[row[2]] = bc
	}
}
