package resources

type ParsedExpense struct {
	Expense    ExpenseResponse `json:"expense"`
	Valid      bool            `json:"valid"`
	Confidence float64         `json:"confidence"`
}
