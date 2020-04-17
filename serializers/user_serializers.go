package serializers

import (
	"ginserver/models"
	"gopkg.in/mgo.v2/bson"
)

// UsersSubsetJSON ..
type UsersSubsetJSON struct {
	ListMeta
	Values []UserSubset `json:"data"`
}

// UserSubset ...
type UserSubset struct {
	UserID   bson.ObjectId `json:"userId,omitempty"`
	Username string        `json:"username"`
	Nickname string        `json:"nickname"`
	Email    string        `json:"email"`
	HasPwd   string        `json:"hasPwd"`
	IsActive string        `json:"isActive"`
	IsVip    string        `json:"isVip"`
	IsStaff  string        `json:"isStaff"`
	IsOper   string        `json:"isOper"`
	IsAdmin  string        `json:"isAdmin"`
}

// NewUserSubset ..
func NewUserSubset(user models.User) UserSubset {
	hasPwd := "false"
	if user.Password != "" {
		hasPwd = "true"
	}
	return UserSubset{
		UserID:   user.UserID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		HasPwd:   hasPwd,
		IsActive: user.IsActive,
		IsVip:    user.IsVip,
	}
}

// NewUsersSubsetJSON ...
func NewUsersSubsetJSON(users []models.User, count int, skip int, total int) UsersSubsetJSON {
	json := UsersSubsetJSON{
		Values: []UserSubset{},
		ListMeta: ListMeta{
			Count: count,
			Skip:  skip,
			Total: total,
		},
	}
	for _, user := range users {
		json.Values = append(json.Values, NewUserSubset(user))
	}
	return json
}

// SerializeUsers ...
func SerializeUsers(users []models.User, params ...interface{}) interface{} {
	count := params[0].(int64)
	skip := params[1].(int64)
	total := params[2].(int)
	userSubsetJSON := NewUsersSubsetJSON(users, int(count), int(skip), total)
	return NewResponse(0, userSubsetJSON, "Success")
}

// SerializeUser ..
func SerializeUser(user models.User) interface{} {
	userSubset := NewUserSubset(user)
	return NewResponse(0, userSubset, "Success")
}
