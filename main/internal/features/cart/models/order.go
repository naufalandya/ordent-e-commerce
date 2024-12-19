package models

import "time"

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	StatusID  int       `json:"status_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User       User                 `json:"user,omitempty"`
	Product    Product              `json:"product,omitempty"`
	Status     OrderStatus          `json:"status,omitempty"`
	OrderItems []TransactionHistory `json:"order_items,omitempty"`
	Payments   []Payment            `json:"payments,omitempty"`
}

type OrderStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TransactionHistory struct {
	ID              int       `json:"id"`
	OrderID         int       `json:"order_id"`
	MidtransOrderID int       `json:"midtrans_order_id,omitempty"`
	Price           float64   `json:"price"`
	Order           Order     `json:"order"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Payment struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ApiResponseCartSuccess struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ApiResponseCartFailed struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
