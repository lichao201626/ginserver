package dao

import (
	"errors"
	"ginserver/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// CreateToken create a token
func CreateToken(mgoClient *mgo.Session, token models.Token) (models.Token, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionToken)
	// add default field
	token.TokenID = bson.NewObjectId()
	// expire after 24 hours
	m, _ := time.ParseDuration("24h")
	token.Expire = time.Now().Add(m).UTC().Format(FormatString)
	token.CreateTime = time.Now().UTC().Format(FormatString)
	token.UpdateTime = time.Now().UTC().Format(FormatString)
	err := collection.Insert(token)
	if err != nil {
		log.Error("Create token failed ", err)
		return models.Token{}, err
	}
	return token, nil
}

// GetTokens ...
func GetTokens(mgoClient *mgo.Session, params ...interface{}) ([]models.Token, error) {
	count, _ := params[0].(int64)
	skip, _ := params[1].(int64)
	query, _ := params[2].(bson.M)

	tokens := []models.Token{}
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionToken)
	err := collection.Find(query).Limit(int(count)).Skip(int(skip)).All(&tokens)
	if err != nil {
		log.Error("Get tokens failed ", err)
		return []models.Token{}, err
	}
	return tokens, nil
}

// GetTokensTotal ...
func GetTokensTotal(mgoClient *mgo.Session, params ...interface{}) (int, error) {
	query, _ := params[0].(bson.M)
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionToken)

	total, err := collection.Find(query).Count()
	if err != nil {
		log.Error("Get tokens total failed ", err)
		return 0, err
	}
	return total, nil
}

// GetTokensByID ...
func GetTokensByID(mgoClient *mgo.Session, tokenID string) (models.Token, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionToken)
	token := models.Token{}

	if bson.IsObjectIdHex(tokenID) {
		err := collection.FindId(bson.ObjectIdHex(tokenID)).One(&token)
		if err != nil {
			log.Error("Get tokens by id failed ", err)
			return token, err
		}
		return token, nil
	}
	return token, errors.New("Not valid bson object id")
}

// UpdateTokensByID ...
func UpdateTokensByID(mgoClient *mgo.Session, params ...interface{}) (models.Token, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionToken)
	tokenID := params[0].(string)
	token := params[1].(models.Token)
	token.UpdateTime = time.Now().UTC().Format(FormatString)

	if bson.IsObjectIdHex(tokenID) {
		err := collection.UpdateId(bson.ObjectIdHex(tokenID), bson.M{
			"$set": token,
		})
		if err != nil {
			log.Error("Update tokens by id failed ", err)
			return token, err
		}
		return token, nil
	}
	return token, errors.New("Not valid bson object id")
}

// RemoveTokens ...
func RemoveTokens(mgoClient *mgo.Session, query bson.M) error {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionToken)
	_, err := collection.RemoveAll(query)
	if err != nil {
		log.Error("Remove tokens failed ", err)
		return err
	}
	return nil
}
