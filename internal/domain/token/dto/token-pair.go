package dto

type TokenPair struct {
	AccessToken      string
	TokenType        string
	ExpiresIn        int
	RefreshToken     string `json:"-"`
	RefreshExpiresIn int    `json:"-"`
}
