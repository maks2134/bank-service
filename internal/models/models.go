package models

import "time"

type BankAccount struct {
	ID            int     `json:"id"`
	AccountType   string  `json:"accounttype"`
	AccountNumber string  `json:"accountnumber"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
	OpenDate      string  `json:"opendate"`
	AccountStatus string  `json:"accountstatus"`
	ClientID      *int    `json:"clientid"`
}

type BankBranch struct {
	ID           int    `json:"id"`
	ServiceZone  string `json:"servicezone"`
	FullAddress  string `json:"full_address"`
	ContactPhone string `json:"contactphone"`
	WorkingDays  string `json:"workingdays"`
	WorkingHours string `json:"workinghours"`
	BranchType   string `json:"branchtype"`
}

type BankStaff struct {
	ID            int     `json:"id"`
	FullName      string  `json:"fullname"`
	Passport      string  `json:"passport"`
	Position      string  `json:"position"`
	HireDate      string  `json:"hiredate"`
	AccessLevel   string  `json:"accesslevel"`
	Qualification *string `json:"qualification"`
	BranchID      int     `json:"branchid"`
}

type Client struct {
	ID         int    `json:"id"`
	Phone      string `json:"phone"`
	Login      string `json:"login"`
	FullName   string `json:"fullname"`
	ClientType string `json:"clienttype"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Passport   string `json:"passport"`
	CreditID   *int   `json:"creditid"`
}

type Credit struct {
	ID            int     `json:"id"`
	Amount        float64 `json:"amount"`
	InterestRate  float64 `json:"interestrate"`
	Purpose       string  `json:"purpose"`
	IssueDate     string  `json:"issuedate"`
	Status        string  `json:"status"`
	RepaymentDate *string `json:"repaymentdate"`
	Currency      string  `json:"currency"`
	RiskLevel     *string `json:"risk_level"`
}

type CreditBankStaff struct {
	CreditID int `json:"creditid"`
	StaffID  int `json:"staffid"`
}

type Transaction struct {
	ID              int       `json:"id"`
	Amount          float64   `json:"amount"`
	OperationDate   time.Time `json:"operationdate"`
	OperationType   string    `json:"operationtype"`
	Purpose         *string   `json:"purpose"`
	OperationStatus string    `json:"operationstatus"`
	Currency        string    `json:"currency"`
	BranchID        int       `json:"branchid"`
	AccountID       int       `json:"accountid"`
}
