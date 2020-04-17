package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	// CollectionOrder is the mongo collection name of order
	CollectionOrder = "orders"
)

// Order is a record of a person's AIP order history
type Order struct {
	OrderID    bson.ObjectId `json:"orderId,omitempty" bson:"_id,omitempty"`
	PlanID     string        `json:"planId" bson:"planId,omitempty"`
	UserID     string        `json:"userId" bson:"userId"`
	Exchange   string        `json:"exchange" bson:"exchange"`
	Price      string        `json:"price" bson:"price"`
	Unit       string        `json:"unit" bson:"unit"`
	Amount     string        `json:"amount" bson:"amount"`
	Fee        string        `json:"fee" bson:"fee"`
	OrderTime  string        `json:"orderTime" bson:"orderTime"`
	CreateTime string        `json:"createTime" bson:"createTime"`
	UpdateTime string        `json:"updateTime" bson:"updateTime"`
}
