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

func GetUserByEmailReturnIDAndRole(email string) (*models.User, error) {
	var user models.User
	var roles []string

	query := `
        SELECT u.id, u.email, u.password, u.full_name, ARRAY_AGG(r.name)
        FROM users u
        LEFT JOIN user_roles ur ON u.id = ur.user_id
        LEFT JOIN roles r ON ur.role_id = r.id
        WHERE u.email = $1
        GROUP BY u.id;
    `

	row := repositories.DB.QueryRow(context.Background(), query, email)

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &roles)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user not found")
		}
		log.Println("Error scanning user with roles:", err)
		return nil, err
	}

	user.Role = roles
	return &user, nil
}

func CreateUser(user *models.User, roleID int) error {
	userQuery := `
		INSERT INTO users (email, full_name, password, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		RETURNING id;
	`

	var newUserID int
	err := repositories.DB.QueryRow(context.Background(), userQuery, user.Email, user.Name, user.Password).Scan(&newUserID)
	if err != nil {
		log.Println("Error inserting user into database:", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = newUserID

	roleQuery := `
		INSERT INTO user_roles (user_id, role_id, created_at)
		VALUES ($1, $2, now());
	`

	_, err = repositories.DB.Exec(context.Background(), roleQuery, newUserID, roleID)
	if err != nil {
		log.Println("Error assigning role to user:", err)
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}
