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

	if err := storage.DB.CreateUser(userToEntityUser(*user)); err != nil {
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

//GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(username string) (User, error) {
	user, err := storage.DB.GetUserByUsername(username)
	if err != nil {
		return User{}, err
	}
	return entityUserToUser(*user), nil
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
