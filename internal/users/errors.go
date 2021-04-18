package users

type WrongUsernameOrPasswordError struct{}

func (m *WrongUsernameOrPasswordError) Error() string {
	return "wrong username or password"
}

type AccessDeniedError struct{}

func (e *AccessDeniedError) Error() string {
	return "access denied"
}

type InvalidUsernameOrPassordError struct{}

func (e *InvalidUsernameOrPassordError) Error() string {
	return "invalid username or password"
}
