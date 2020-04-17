package dao

import (
	"errors"
	"ginserver/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// CreatePincode in mongo
func CreatePincode(mgoClient *mgo.Session, pincode models.Pincode) (models.Pincode, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPincode)
	// add default field
	pincode.PincodeID = bson.NewObjectId()
	// expire after 15 minute
	m, _ := time.ParseDuration("15m")
	pincode.Tried = 0
	pincode.Expire = time.Now().Add(m).UTC().Format(FormatString)
	pincode.CreateTime = time.Now().UTC().Format(FormatString)
	pincode.UpdateTime = time.Now().UTC().Format(FormatString)
	err := collection.Insert(pincode)
	if err != nil {
		log.Error("Create pincode failed ", err)
		return models.Pincode{}, err
	}
	return pincode, nil
}

// GetPincodes ...
func GetPincodes(mgoClient *mgo.Session, params ...interface{}) ([]models.Pincode, error) {
	count, _ := params[0].(int64)
	skip, _ := params[1].(int64)
	query, _ := params[2].(bson.M)

	pincodes := []models.Pincode{}
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPincode)
	err := collection.Find(query).Sort("-createTime").Skip(int(skip)).Limit(int(count)).All(&pincodes)
	if err != nil {
		log.Error("Get pincodes failed ", err)
		return []models.Pincode{}, err
	}
	return pincodes, nil
}

// GetPincodesTotal ...
func GetPincodesTotal(mgoClient *mgo.Session, params ...interface{}) (int, error) {
	query, _ := params[0].(bson.M)
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPincode)

	total, err := collection.Find(query).Count()
	if err != nil {
		log.Error("Get pincodes total failed ", err)
		return 0, err
	}
	return total, nil
}

// GetPincodesByID ...
func GetPincodesByID(mgoClient *mgo.Session, pincodeID string) (models.Pincode, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPincode)
	pincode := models.Pincode{}

	if bson.IsObjectIdHex(pincodeID) {
		err := collection.FindId(bson.ObjectIdHex(pincodeID)).One(&pincode)
		if err != nil {
			log.Error("Get pincode by id failed ", err)
			return pincode, err
		}
		return pincode, nil
	}
	return pincode, errors.New("Not valid bson object id")
}

// UpdatePincodesByID ...
func UpdatePincodesByID(mgoClient *mgo.Session, params ...interface{}) (models.Pincode, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPincode)
	pincodeID := params[0].(string)
	pincode := params[1].(models.Pincode)
	pincode.UpdateTime = time.Now().UTC().Format(FormatString)

	if bson.IsObjectIdHex(pincodeID) {
		err := collection.UpdateId(bson.ObjectIdHex(pincodeID), bson.M{
			"$set": pincode,
		})
		if err != nil {
			log.Error("Update pincode by id failed ", err)
			return pincode, err
		}
		return pincode, nil
	}
	return pincode, errors.New("not valid bson object id")
}

// RemovePincodes ...
func RemovePincodes(mgoClient *mgo.Session, query bson.M) error {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPincode)
	_, err := collection.RemoveAll(query)
	if err != nil {
		log.Error("Remove pincode failed ", err)
		return err
	}
	return nil
}
