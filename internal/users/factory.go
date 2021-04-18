package users

import "github.com/Shelex/split-specs/entities"

func userToEntityUser(user User) entities.User {
	return entities.User{
		Username: user.Username,
		Password: user.Password,
		ID:       user.ID,
	}
}

func entityUserToUser(user entities.User) User {
	return User{
		Username: user.Username,
		Password: user.Password,
		ID:       user.ID,
	}
}
