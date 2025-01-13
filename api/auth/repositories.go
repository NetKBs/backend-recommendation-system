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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var email string
	if err := session.Query(`SELECT email FROM user_by_email WHERE email = ?`, user.Email).Scan(&email); err != nil {
		if err != gocql.ErrNotFound {
			return err
		}
	}
	if email != "" {
		return fmt.Errorf("email already exists")
	}

	if err := session.Query(`INSERT INTO user_by_email (email, user_id, name, password) VALUES (?, ?, ?, ?)`, user.Email, newUUID, user.Name, string(hashedPassword)).Exec(); err != nil {
		return err
	}
	if err := session.Query(`INSERT INTO user_by_id (user_id, name, email, password) VALUES (?, ?, ?, ?)`,
		newUUID, user.Name, user.Email, string(hashedPassword)).Exec(); err != nil {
		return err
	}

	return nil
}

func LoginRepository(user schema.UserLogin) error {
	session := config.SESSION
	var existingEmail string
	var hashedPassword string

	if err := session.Query(`SELECT email, password FROM user_by_email WHERE email = ? LIMIT 1`, user.Email).Scan(&existingEmail, &hashedPassword); err != nil {
		if err == gocql.ErrNotFound {
			return fmt.Errorf("email does not exist")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		return fmt.Errorf("wrong password")
	}

	return nil
}
