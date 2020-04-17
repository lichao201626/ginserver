package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	// CollectionPlan is the mongo collection name of plan
	CollectionPlan = "plans"
)

// Plan is a record of a person's AIP
type Plan struct {
	PlanID     bson.ObjectId `json:"planId,omitempty" bson:"_id,omitempty"`
	UserID     string        `json:"userId" bson:"userId"`
	PlanName   string        `json:"planName" bson:"planName"`
	Budget     string        `json:"budget" bson:"budget"`
	Symbol     string        `json:"symbol" bson:"symbol"`
	Frequency  string        `json:"frequency" bson:"frequency"`
	SubTime    string        `json:"subTime" bson:"subTime"`
	APIKeyID   string        `json:"apiKeyId" bson:"apiKeyId"`
	APIKeyName string        `json:"apiKeyName" bson:"apiKeyName"`
	IsActive   string        `json:"isActive" bson:"isActive"`
	CreateTime string        `json:"createTime" bson:"createTime"`
	UpdateTime string        `json:"updateTime" bson:"updateTime"`
}
