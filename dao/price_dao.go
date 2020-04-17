package dao

import (
	"ginserver/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetPriceData ...
func GetPriceData(mgoClient *mgo.Session, params ...interface{}) []models.Price {
	query, _ := params[0].(bson.M)

	prices := []models.Price{}
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPrice)
	err := collection.Find(query).All(&prices)
	if err != nil {
		log.Error("Get prices error ", err)
	}
	return prices
}
