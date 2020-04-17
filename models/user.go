package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	// CollectionUser is the mongo collection name of User
	CollectionUser = "users"
)

// User is a record for a person who can login to the system
type User struct {
	UserID     bson.ObjectId `json:"userId,omitempty" bson:"_id,omitempty"`
	Username   string        `json:"username" bson:"username"`
	Nickname   string        `json:"nickname" bson:"nickname"`
	Email      string        `json:"email" form:"email" binding:"required" bson:"email"`
	Password   string        `json:"password" form:"password" binding:"required" bson:"password"`
	IsActive   string        `json:"isActive" bson:"isActive"`
	IsVip      string        `json:"isVip" bson:"isVip"`
	IsStaff    string        `json:"isStaff" bson:"isStaff"`
	IsOper     string        `json:"isOper" bson:"isOper"`
	IsAdmin    string        `json:"isAdmin" bson:"isAdmin"`
	CreateTime string        `json:"createTime" bson:"createTime"`
	UpdateTime string        `json:"updateTime" bson:"updateTime"`
}

// NewUser instatiates a new user
func NewUser() User {
	return User{}
}
