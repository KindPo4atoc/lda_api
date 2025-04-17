package entity

// структура для хранения ОДНОЙ ЗАПИСИ из бд
type UserData struct {
	LoanId      int    `json:"loan_id"`
	LoanTerm    int    `json:"loan_term"`
	IncomeAnnum int64  `json:"income_annum"`
	LoanAmount  int64  `json:"loan_amount"`
	CibilScore  int    `json:"cibil_score"`
	LoanStatus  string `json:"loan_status"`
}
