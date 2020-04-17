package dao

import (
	"errors"
	"ginserver/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// CreateOrder an order
func CreateOrder(mgoClient *mgo.Session, order models.Order) (models.Order, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionOrder)
	// add default field
	order.OrderID = bson.NewObjectId()
	order.CreateTime = time.Now().UTC().Format(FormatString)
	order.UpdateTime = time.Now().UTC().Format(FormatString)

	err := collection.Insert(order)
	if err != nil {
		log.Error("Create order error ", err)
		return models.Order{}, err
	}
	return order, nil
}

// GetOrders ...
func GetOrders(mgoClient *mgo.Session, params ...interface{}) ([]models.Order, error) {
	count, _ := params[0].(int64)
	skip, _ := params[1].(int64)
	query, _ := params[2].(bson.M)

	orders := []models.Order{}
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionOrder)
	err := collection.Find(query).Sort("-createTime").Skip(int(skip)).Limit(int(count)).All(&orders)
	if err != nil {
		log.Error("Get orders error ", err)
		return []models.Order{}, err
	}
	return orders, nil
}

// GetOrdersTotal ...
func GetOrdersTotal(mgoClient *mgo.Session, params ...interface{}) (int, error) {
	query, _ := params[0].(bson.M)
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionOrder)

	total, err := collection.Find(query).Count()
	if err != nil {
		log.Error("Get orders total error ", err)
		return 0, err
	}
	return total, nil
}

// GetOrdersByID ...
func GetOrdersByID(mgoClient *mgo.Session, orderID string) (models.Order, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionOrder)
	order := models.Order{}

	if bson.IsObjectIdHex(orderID) {
		err := collection.FindId(bson.ObjectIdHex(orderID)).One(&order)
		if err != nil {
			log.Error("Get orders by id error ", err)
			return order, err
		}
		return order, nil
	}
	return order, errors.New("not valid bson object id")
}

// UpdateOrdersByID ...
func UpdateOrdersByID(mgoClient *mgo.Session, params ...interface{}) (models.Order, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionOrder)
	orderID := params[0].(string)
	order := params[1].(models.Order)
	order.UpdateTime = time.Now().UTC().Format(FormatString)

	if bson.IsObjectIdHex(orderID) {
		err := collection.UpdateId(bson.ObjectIdHex(orderID), bson.M{
			"$set": order,
		})
		if err != nil {
			log.Error("Update orders by id error ", err)
			return order, err
		}
		return order, nil
	}
	return order, errors.New("not valid bson object id")
}
