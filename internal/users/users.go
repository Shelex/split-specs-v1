package users

import (
	"github.com/Shelex/split-specs/storage"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
	Password string
}

func (user *User) Create() error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	if err := storage.DB.CreateUser(UserToEntityUser(*user)); err != nil {
		return err
	}
	return nil
}

func (user *User) Authenticate() bool {
	dbUser, err := storage.DB.GetUserByUsername(user.Username)
	if err != nil {
		return false
	}
	return CheckPasswordHash(user.Password, dbUser.Password)
}

func (user *User) Exist() bool {
	if _, err := storage.DB.GetUserByUsername(user.Username); err != nil {
		return false
	}
	return true
}

func (user *User) ChangePassword(password string, newPassword string) error {
	dbUser, err := storage.DB.GetUserByUsername(user.Username)
	if err != nil {
		return &AccessDeniedError{}
	}
	if match := CheckPasswordHash(password, dbUser.Password); !match {
		return &AccessDeniedError{}
	}
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	return storage.DB.UpdatePassword(user.ID, hashedPassword)
}

//HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
