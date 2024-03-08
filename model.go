package models

// Untuk User
type Users struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type UsersResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    []Users `json:"data"`
}

// Untuk Products
type Products struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}
type ProductsResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []Products `json:"data"`
}

// Untuk Transactions
type Transactions struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
type TransactionResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    Transactions `json:"data"`
}
type TransactionsResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    []Transactions `json:"data"`
}

// Untuk Detail Transactions
type DetailTransactions struct {
	Users        Users        `json:"users"`
	Products     Products     `json:"products"`
	Transactions Transactions `json:"transactions"`
}
type DetailTransactionsResponse struct {
	Status  int                  `json:"status"`
	Message string               `json:"message"`
	Data    []DetailTransactions `json:"data"`
}

// Error response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Success response
type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
