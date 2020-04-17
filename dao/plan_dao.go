package dao

import (
	"errors"
	"ginserver/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// CreatePlan create a plan
func CreatePlan(mgoClient *mgo.Session, plan models.Plan) (models.Plan, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPlan)
	// add default field
	plan.PlanID = bson.NewObjectId()
	plan.CreateTime = time.Now().UTC().Format(FormatString)
	plan.UpdateTime = time.Now().UTC().Format(FormatString)

	err := collection.Insert(plan)
	if err != nil {
		log.Error("Create plan error ", err)
		return models.Plan{}, err
	}
	return plan, nil
}

// GetPlans ...
func GetPlans(mgoClient *mgo.Session, params ...interface{}) ([]models.Plan, error) {
	count, _ := params[0].(int64)
	skip, _ := params[1].(int64)
	query, _ := params[2].(bson.M)

	plans := []models.Plan{}
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPlan)
	err := collection.Find(query).Sort("-createTime").Skip(int(skip)).Limit(int(count)).All(&plans)
	if err != nil {
		log.Error("Get plans error ", err)
		return []models.Plan{}, err
	}
	return plans, nil
}

// GetPlansTotal ...
func GetPlansTotal(mgoClient *mgo.Session, params ...interface{}) (int, error) {
	query, _ := params[0].(bson.M)
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPlan)

	total, err := collection.Find(query).Count()
	if err != nil {
		log.Error("Get Plans total error ", err)
		return 0, err
	}
	return total, nil
}

// GetPlansByID ...
func GetPlansByID(mgoClient *mgo.Session, planID string) (models.Plan, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPlan)
	plan := models.Plan{}

	if bson.IsObjectIdHex(planID) {
		err := collection.FindId(bson.ObjectIdHex(planID)).One(&plan)
		if err != nil {
			log.Error("Get plans by id error ", err)
			return plan, err
		}
		return plan, nil
	}
	return plan, errors.New("not valid bson object id")
}

// UpdatePlansByID ...
func UpdatePlansByID(mgoClient *mgo.Session, params ...interface{}) (models.Plan, error) {
	collection := mgoClient.DB("quanta_lab_aip").C(models.CollectionPlan)
	planID := params[0].(string)
	plan := params[1].(models.Plan)
	plan.UpdateTime = time.Now().UTC().Format(FormatString)

	if bson.IsObjectIdHex(planID) {
		err := collection.UpdateId(bson.ObjectIdHex(planID), bson.M{
			"$set": plan,
		})
		if err != nil {
			log.Error("Update Plans by id error ", err)
			return plan, err
		}
		return plan, nil
	}
	return plan, errors.New("not valid bson object id")
}
