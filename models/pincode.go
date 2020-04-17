package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	// CollectionPincode is the mongo collection name of Pincode
	CollectionPincode = "Pincodes"
)

// Pincode is used to verify email
type Pincode struct {
	PincodeID  bson.ObjectId `json:"pincodeId,omitempty" bson:"_id,omitempty"`
	PinCode    string        `json:"pincode" bson:"pincode"`
	Expire     string        `json:"expire" bson:"expire"`
	UserID     string        `json:"userId" bson:"userId"`
	Tried      int           `json:"tried" bson:"tried"`
	CreateTime string        `json:"createTime" bson:"createTime"`
	UpdateTime string        `json:"updateTime" bson:"updateTime"`
}

// NewPincode instatiates a new Pincode
func NewPincode() Pincode {
	return Pincode{}
}
