package auth

type AuthDefault struct {
	Token string
}

func NewAuthDefault(token string) *AuthDefault {
	return &AuthDefault{
		Token: token,
	}
}

// Auth checks if the token is valid
func (a *AuthDefault) Auth(token string) (err error) {
	if a.Token == token {
		return nil
	}

	return ErrAuthTokenInvalid
}

// GetToken returns the token
func (a *AuthDefault) GetToken() string {
	return a.Token
}
