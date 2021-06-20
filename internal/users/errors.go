package users

type InvalidEmailFormat struct{}

func (m *InvalidEmailFormat) Error() string {
	return "wrong email format"
}

type WrongEmailOrPasswordError struct{}

func (m *WrongEmailOrPasswordError) Error() string {
	return "wrong email or password"
}

type AccessDeniedError struct{}

func (e *AccessDeniedError) Error() string {
	return "access denied"
}

type InvalidEmailOrPassordError struct{}

func (e *InvalidEmailOrPassordError) Error() string {
	return "invalid email or password"
}
