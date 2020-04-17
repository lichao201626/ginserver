package resources

import (
	"encoding/json"
	"fmt"
	"ginserver/dao"
	"ginserver/models"
	"ginserver/serializers"
	//"ginserver/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
)

// PlanResource ...
type PlanResource struct {
}

// NewPlanResource ...
func NewPlanResource(e *gin.Engine) {
	u := PlanResource{}
	// Setup Routes
	e.GET("/get/plans", u.getPlans)
	e.GET("/get/plans/:id", u.getPlansByID)
	e.POST("/post/plans", u.postPlans)
	e.PATCH("/patch/plans/:id", u.patchPlansByID)
}

type planJSON struct {
	UserID string `json:"userId"`
	//plan
	PlanName  string `json:"planName"`
	Budget    string `json:"budget"`
	Symbol    string `json:"symbol"`
	Frequency string `json:"frequency"`
	SubTime   string `json:"subTime"`
	IsActive  string `json:"isActive"`
	//key
	Exchange   string `json:"exchange"`
	Key        string `json:"key"`
	Secret     string `json:"secret"`
	Phrase     string `json:"phrase"`
	APIKeyName string `json:"apiKeyName"`
}

// ValidateRes .. when validate the plan and key
type ValidateRes struct {
	IsPlanValid    string `json:"isPlanValid"`
	IsKeyValid     string `json:"isKeyValid"`
	IsBudgetEnough string `json:"isBudgetEnough"`
}

func (r *PlanResource) postPlans(c *gin.Context) {
	plan := planJSON{}
	c.MustBindWith(&plan, binding.JSON)
	currentUser := c.MustGet("currentUser").(models.User)
	// todo may need user role verify
	//if currentUser.IsActive or IsVip
	plan.UserID = currentUser.UserID.Hex()

	isBodyValid := validatePostPlan(c, plan)
	if !isBodyValid {
		return
	}

	res, err := postValidate(plan)
	if err != nil || res.Body == nil {
		c.JSON(500, serializers.SerializeError(50001, "Validate plan error"))
		return
	}

	log.Info(res)
	log.Info(res.Body)
	log.Info(res.Code)
	validateJSON := res.Body.(map[string]interface{})
	/* 	validateJSON := ValidateRes{
		IsPlanValid:    "true",
		IsKeyValid:     "true",
		IsBudgetEnough: "true",
	} */
	if validateJSON["isKeyValid"].(string) != "true" {
		c.JSON(200, serializers.SerializeError(20001, "Key is invalid"))
		return
	}
	if validateJSON["isPlanValid"].(string) != "true" {
		c.JSON(200, serializers.SerializeError(20002, "Plan is invalid"))
		return
	}
	if validateJSON["isBudgetEnough"].(string) != "true" {
		c.JSON(200, serializers.SerializeError(20003, "Budget is not enough"))
		return
	}

	// save mgo
	// save key , if key exist,

	// save plan

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	query := bson.M{
		"userId":     plan.UserID,
		"apiKeyName": plan.APIKeyName,
	}
	keys, err := dao.GetKeys(dbSession, 20, 0, query)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Mongo get keys failed"))
		return
	}

	keyObj := models.Key{
		UserID:     plan.UserID,
		APIKeyName: plan.APIKeyName,
		Key:        plan.Key,
		Secret:     plan.Secret,
		Phrase:     plan.Phrase,
	}
	// encrypt the key
	encryptErr := keyObj.Encrypt()
	if encryptErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Encrypt the key failed"))
		return
	}
	// save the key
	if len(keys) > 0 {
		keyObj = keys[0]
	} else {
		keyObj, err = dao.CreateKey(dbSession, keyObj)
		if err != nil {
			c.JSON(500, serializers.SerializeError(50002, "Mongo create key failed"))
			return
		}
	}
	// save the plan
	planObj := models.Plan{
		APIKeyID:  keyObj.APIKeyID.Hex(),
		UserID:    plan.UserID,
		PlanName:  plan.PlanName,
		Budget:    plan.Budget,
		Symbol:    plan.Symbol,
		Frequency: plan.Frequency,
		SubTime:   plan.SubTime,
		IsActive:  plan.IsActive,
	}
	planObj, err = dao.CreatePlan(dbSession, planObj)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50003, "Mongo create plan failed"))
		return
	}
	log.Info("New Plan created ", planObj)
	// todo may email notify users when plan created success
	c.JSON(200, serializers.SerializePlan(planObj))
}

func (r *PlanResource) getPlans(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	/* 	if currentUser.IsAdmin != "true" {
		c.JSON(401, serializers.SerializeError(40102, "Can not get info of other Plan"))
		return
	} */
	skip, _ := getSkip(c)
	count, _ := getCount(c)
	createBefore := c.Query("createBefore")
	createAfter := c.Query("createAfter")
	Planname := c.Query("Planname")
	log.Info("Planname:", Planname)

	query := bson.M{}
	query["userId"] = currentUser.UserID.Hex()
	if createAfter != "" && createBefore == "" {
		query["createTime"] = bson.M{
			"$gte": createAfter,
		}
	}
	if createBefore != "" && createAfter == "" {
		query["createTime"] = bson.M{
			"$lte": createBefore,
		}
	}
	if createBefore != "" && createAfter != "" {
		query["createTime"] = bson.M{
			"$gte": createAfter,
			"$lte": createBefore,
		}
	}
	if Planname != "" {
		query["Planname"] = Planname
	}

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	total, err := dao.GetPlansTotal(dbSession, query)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Mongo get plans total failed"))
		return
	}
	plans, getErr := dao.GetPlans(dbSession, count, skip, query)
	if getErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo get Plans failed"))
		return
	}
	log.Info("Get plans success ", total)
	c.JSON(200, serializers.SerializePlans(plans, count, skip, total))
}

func (r *PlanResource) getPlansByID(c *gin.Context) {
	planID := c.Param("id")
	// currentUserID := c.MustGet("currentUserID").(string)
	currentUser := c.MustGet("currentUser").(models.User)

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	plan, err := dao.GetPlansByID(dbSession, planID)
	if err != nil {
		c.JSON(200, serializers.SerializeError(20001, "Not found by the plan id"))
		return
	}
	if plan.UserID != currentUser.UserID.Hex() {
		c.JSON(401, serializers.SerializeError(40102, "Can not get info of other Plan"))
		return
	}
	log.Info("Get Plans by id success ", planID)
	c.JSON(200, serializers.SerializePlan(plan))
}

func (r *PlanResource) patchPlansByID(c *gin.Context) {
	planID := c.Param("id")
	currentUser, _ := c.MustGet("currentUser").(models.User)

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	plan, err := dao.GetPlansByID(dbSession, planID)
	if err != nil {
		c.JSON(200, serializers.SerializeError(20001, "Not found by the Plan id"))
		return
	}
	if plan.UserID != currentUser.UserID.Hex() {
		c.JSON(401, serializers.SerializeError(40101, "Can not change info of other Plan"))
		return
	}

	body := getPatchBody(c)
	if len(body) == 0 {
		c.JSON(400, serializers.SerializeError(40001, "Empty patch body, or body invalid format"))
		return
	}

	// PlanName := Plan.Planname
	// planReq := planJSON{}
	applyPatchPlanBody(&plan, body)

	log.Info("plan:", plan)
	pJSON := planJSON{
		UserID:    plan.UserID,
		PlanName:  plan.PlanName,
		Budget:    plan.Budget,
		Symbol:    plan.Symbol,
		Frequency: plan.Frequency,
		SubTime:   plan.SubTime,
		IsActive:  plan.IsActive,
	}
	res, err := postValidate(pJSON)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Validate plan error"))
		return
	}
	log.Info(res)
	validateJSON := res.Body.(map[string]interface{})
	/* 	validateJSON := ValidateRes{
		IsPlanValid:    "true",
		IsKeyValid:     "true",
		IsBudgetEnough: "true",
	} */
	if validateJSON["isKeyValid"].(string) != "true" {
		c.JSON(200, serializers.SerializeError(20001, "Key is invalid"))
		return
	}
	if validateJSON["isPlanValid"].(string) != "true" {
		c.JSON(200, serializers.SerializeError(20002, "Plan is invalid"))
		return
	}
	if validateJSON["isBudgetEnough"].(string) != "true" {
		c.JSON(200, serializers.SerializeError(20003, "Budget is not enough"))
		return
	}
	// save mgo
	// save key , if key exist,
	// save plan
	plan, err = dao.UpdatePlansByID(dbSession, planID, plan)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo update Plans by id failed"))
		return
	}
	log.Info("Patch Plan by id success ", planID)
	c.JSON(200, serializers.SerializePlan(plan))
}

func applyPatchPlanBody(plan *models.Plan, body []PatchJSON) {
	for _, v := range body {
		if v.Op != "update" {
			continue
		}

		switch v.Path {
		case "/planName":
			plan.PlanName = v.Value
		case "/budget":
			plan.Budget = v.Value
		case "/symbol":
			plan.Symbol = v.Value
		case "/frequency":
			plan.Frequency = v.Value
		case "/subTime":
			plan.SubTime = v.Value
		case "/isActive":
			plan.IsActive = v.Value
		case "/apiKeyId":
			plan.APIKeyID = v.Value
		default:
			continue
		}
	}
}

func validatePostPlan(c *gin.Context, plan planJSON) bool {
	log.Info(plan)
	if plan.PlanName == "" {
		c.JSON(400, serializers.SerializeError(40001, "Must supply planName"))
		return false
	}

	if plan.Budget == "" {
		c.JSON(400, serializers.SerializeError(40002, "Must supply budget"))
		return false
	}

	if plan.Symbol == "" {
		c.JSON(400, serializers.SerializeError(40003, "Must supply symbol"))
		return false
	}

	if plan.Frequency == "" {
		c.JSON(400, serializers.SerializeError(40004, "Must supply frequency"))
		return false
	}

	if plan.SubTime == "" {
		c.JSON(400, serializers.SerializeError(40005, "Must supply subTime"))
		return false
	}

	if plan.IsActive == "" {
		c.JSON(400, serializers.SerializeError(40006, "Must supply isActive"))
		return false
	}

	if plan.Exchange == "" {
		c.JSON(400, serializers.SerializeError(40007, "Must supply exchange"))
		return false
	}

	if plan.Key == "" {
		c.JSON(400, serializers.SerializeError(40008, "Must supply key"))
		return false
	}

	if plan.Secret == "" {
		c.JSON(400, serializers.SerializeError(40009, "Must supply secret"))
		return false
	}

	/* 	if plan.Phrase == "" {
		c.JSON(400, serializers.SerializeError(40010, "Must supply phrase"))
		return false
	} */

	if plan.APIKeyName == "" {
		c.JSON(400, serializers.SerializeError(40011, "Must supply api key name"))
		return false
	}

	return true
}

func postValidate(plan planJSON) (resData, error) {
	//url := "http://52.76.225.208:8687/validatePlan"
	// url := "http://10.200.10.17:8080/validatePlan"
	url := "http://localhost:8889/post/validate"
	planStr, err := json.Marshal(plan)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(planStr)))
	if err != nil {
		// panic(err)

		return resData{}, err
	}
	// key, sign := encryptReq("/validatePlan", "post", plan)
	key, sign := encryptReq("/post/validate", "post", plan)
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("BTAI-ACCESS-KEY", key)
	req.Header.Set("BTAI-ACCESS-SIGN", sign)
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "Keep-Alive")
	// req.Body.Read

	client := &http.Client{}
	fmt.Println("req", req)
	resp, err := client.Do(req)
	if err != nil {
		// panic(err)
		return resData{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("URL:>", body)
	var data resData
	json.Unmarshal([]byte(string(body)), &data)
	fmt.Println(data)
	// var data []priceData
	// json.Unmarshal([]byte(string(body)), &data)
	fmt.Println(string(body))
	return data, nil
}
