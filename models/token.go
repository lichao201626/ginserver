package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	// CollectionToken is the mango collection name of Token
	CollectionToken = "Tokens"
)

// Token is a record for a person who can login to the system
type Token struct {
	TokenID    bson.ObjectId `json:"tokenId,omitempty" bson:"_id,omitempty"`
	AuthToken  string        `json:"authToken" bson:"authToken"`
	Expire     string        `json:"expire" bson:"expire"`
	UserID     string        `json:"userId" bson:"userId"`
	Type       string        `json:"type" bson:"type"`
	CreateTime string        `json:"createTime" bson:"createTime"`
	UpdateTime string        `json:"updateTime" bson:"updateTime"`
}

// NewToken instatiates a new Token
func NewToken() Token {
	return Token{}
}
