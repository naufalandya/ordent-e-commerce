package models

type User struct {
	ID       int      `json:"id"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Name     string   `json:"name"`
	Role     []string `json:"roles"`
}

type ApiResponseSuccess struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ApiResponseFailed struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
