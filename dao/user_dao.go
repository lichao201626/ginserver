package dao

import (
	"errors"
	"ginserver/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// FormatString is const for date format
const FormatString string = "2006-01-02T15:04:05.000Z"

// CreateUser an user
func CreateUser(mgoClient *mgo.Session, user models.User) (models.User, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionUser)
	// add default field
	user.UserID = bson.NewObjectId()
	user.Password = ""
	user.CreateTime = time.Now().UTC().Format(FormatString)
	user.UpdateTime = time.Now().UTC().Format(FormatString)
	user.IsActive = "false"
	user.IsAdmin = "false"
	user.IsActive = "false"
	user.IsOper = "false"
	user.IsStaff = "false"
	user.IsVip = "false"

	err := collection.Insert(user)
	if err != nil {
		log.Error("Create user error ", err)
		return models.User{}, err
	}
	return user, nil
}

// GetUsers ...
func GetUsers(mgoClient *mgo.Session, params ...interface{}) ([]models.User, error) {
	count, _ := params[0].(int64)
	skip, _ := params[1].(int64)
	query, _ := params[2].(bson.M)

	users := []models.User{}
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionUser)
	err := collection.Find(query).Sort("-createTime").Skip(int(skip)).Limit(int(count)).All(&users)
	if err != nil {
		log.Error("Get users error ", err)
		return []models.User{}, err
	}
	return users, nil
}

// GetUsersTotal ...
func GetUsersTotal(mgoClient *mgo.Session, params ...interface{}) (int, error) {
	query, _ := params[0].(bson.M)
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionUser)

	total, err := collection.Find(query).Count()
	if err != nil {
		log.Error("Get users total error ", err)
		return 0, err
	}
	return total, nil
}

// GetUsersByID ...
func GetUsersByID(mgoClient *mgo.Session, userID string) (models.User, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionUser)
	user := models.User{}

	if bson.IsObjectIdHex(userID) {
		err := collection.FindId(bson.ObjectIdHex(userID)).One(&user)
		if err != nil {
			log.Error("Get users by id error ", err)
			return user, err
		}
		return user, nil
	}
	return user, errors.New("not valid bson object id")
}

// UpdateUsersByID ...
func UpdateUsersByID(mgoClient *mgo.Session, params ...interface{}) (models.User, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionUser)
	userID := params[0].(string)
	user := params[1].(models.User)
	user.UpdateTime = time.Now().UTC().Format(FormatString)

	if bson.IsObjectIdHex(userID) {
		err := collection.UpdateId(bson.ObjectIdHex(userID), bson.M{
			"$set": user,
		})
		if err != nil {
			log.Error("Update users by id error ", err)
			return user, err
		}
		return user, nil
	}
	return user, errors.New("not valid bson object id")
}
