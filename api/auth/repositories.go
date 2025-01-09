package auth

import (
	"example/api/schema"
	"example/config"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRepository(user schema.UserRegister) error {
	session := config.SESSION
	newUUID := gocql.UUID(uuid.New())

	var email string
	var userID gocql.UUID
	applied, err := session.Query(`INSERT INTO uniqueness_emails (email, user_id) VALUES (?, ?) IF NOT EXISTS`, user.Email, newUUID).ScanCAS(&email, &userID)
	if err != nil {
		return err
	}
	if !applied {
		return fmt.Errorf("email already exists")
	}

	if err := session.Query(`INSERT INTO users (user_id, name, email, password) VALUES (?, ?, ?, ?)`,
		newUUID, user.Name, user.Email, user.Password).Exec(); err != nil {
		return err
	}

	return nil
}

func LoginRepository(user schema.UserLogin) error {
	session := config.SESSION
	var existingEmail string
	var existingUserID gocql.UUID

	if err := session.Query(`SELECT email, user_id FROM uniqueness_emails WHERE email = ? LIMIT 1`, user.Email).Scan(&existingEmail, &existingUserID); err != nil {
		if err == gocql.ErrNotFound {
			return fmt.Errorf("email does not exist")
		}
		return err
	}

	var hashedPassword string
	if err := session.Query(`SELECT password FROM users WHERE user_id = ?`, existingUserID).Scan(&hashedPassword); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		return fmt.Errorf("wrong password")
	}

	return nil
}
