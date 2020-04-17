package serializers

import "ginserver/models"

// TokenSubset ...
type TokenSubset struct {
	UserID    string `json:"userId"`
	AuthToken string `json:"authToken"`
}

// NewTokenSubset ..
func NewTokenSubset(token models.Token) TokenSubset {
	return TokenSubset{
		UserID:    token.UserID,
		AuthToken: token.AuthToken,
	}
}

// SerializeToken ..
func SerializeToken(token models.Token) interface{} {
	tokenSubset := NewTokenSubset(token)
	return NewResponse(0, tokenSubset, "Success")
}
