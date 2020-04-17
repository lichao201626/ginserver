package dao

import (
	"errors"
	"ginserver/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// CreateKey create a key
func CreateKey(mgoClient *mgo.Session, key models.Key) (models.Key, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionKey)
	// add default field
	key.APIKeyID = bson.NewObjectId()
	key.CreateTime = time.Now().UTC().Format(FormatString)
	key.UpdateTime = time.Now().UTC().Format(FormatString)

	err := collection.Insert(key)
	if err != nil {
		log.Error("Create key error ", err)
		return models.Key{}, err
	}
	return key, nil
}

// GetKeys ...
func GetKeys(mgoClient *mgo.Session, params ...interface{}) ([]models.Key, error) {
	count, _ := params[0].(int64)
	skip, _ := params[1].(int64)
	query, _ := params[2].(bson.M)

	keys := []models.Key{}
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionKey)
	err := collection.Find(query).Sort("-createTime").Skip(int(skip)).Limit(int(count)).All(&keys)
	if err != nil {
		log.Error("Get keys error ", err)
		return []models.Key{}, err
	}
	return keys, nil
}

// GetKeysTotal ...
func GetKeysTotal(mgoClient *mgo.Session, params ...interface{}) (int, error) {
	query, _ := params[0].(bson.M)
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionKey)

	total, err := collection.Find(query).Count()
	if err != nil {
		log.Error("Get keys total error ", err)
		return 0, err
	}
	return total, nil
}

// GetKeysByID ...
func GetKeysByID(mgoClient *mgo.Session, apiKeyID string) (models.Key, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionKey)
	key := models.Key{}

	if bson.IsObjectIdHex(apiKeyID) {
		err := collection.FindId(bson.ObjectIdHex(apiKeyID)).One(&key)
		if err != nil {
			log.Error("Get keys by id error ", err)
			return key, err
		}
		return key, nil
	}
	return key, errors.New("not valid bson object id")
}

// UpdateKeysByID ...
func UpdateKeysByID(mgoClient *mgo.Session, params ...interface{}) (models.Key, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionKey)
	apiKeyID := params[0].(string)
	key := params[1].(models.Key)
	key.UpdateTime = time.Now().UTC().Format(FormatString)

	if bson.IsObjectIdHex(apiKeyID) {
		err := collection.UpdateId(bson.ObjectIdHex(apiKeyID), bson.M{
			"$set": key,
		})
		if err != nil {
			log.Error("Update keys by id error ", err)
			return key, err
		}
		return key, nil
	}
	return key, errors.New("not valid bson object id")
}

// RemoveKeys ...
func RemoveKeys(mgoClient *mgo.Session, query bson.M) error {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPincode)
	_, err := collection.RemoveAll(query)
	if err != nil {
		log.Error("Remove keys failed ", err)
		return err
	}
	return nil
}
