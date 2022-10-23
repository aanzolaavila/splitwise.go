package resources

type ParsedExpense struct {
	Expense    Expense `json:"expense"`
	Valid      bool    `json:"valid"`
	Confidence float64 `json:"confidence"`
	Error      string  `json:"error"`
}
