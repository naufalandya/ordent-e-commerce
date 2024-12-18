package services

import (
	"commerce/internal/features/auth/models"
	"commerce/internal/repositories"
	"context"
	"fmt"
	"log"
)

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	row := repositories.DB.QueryRow(context.Background(), "SELECT id, email, password, full_name FROM users WHERE email = $1", email)

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user not found")
		}
		log.Println("Error scanning user:", err)
		return nil, err
	}

	return &user, nil
}
