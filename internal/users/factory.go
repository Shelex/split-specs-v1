package users

import "github.com/Shelex/split-specs/entities"

func UserToEntityUser(user User) entities.User {
	return entities.User{
		Username: user.Username,
		Password: user.Password,
		ID:       user.ID,
	}
}

func EntityUserToUser(user entities.User) User {
	return User{
		Username: user.Username,
		Password: user.Password,
		ID:       user.ID,
	}
}
