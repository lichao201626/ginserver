package resources

import (
	"ginserver/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	// "gopkg.in/mgo.v2/bson"
	"reflect"
	"regexp"
	"strconv"
	// "strings"
)

// EmailPattern is pattern to test if a email address is valid
const EmailPattern = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`

// UsernamePattern is pattern to test if a username is valid
// const UsernamePattern = `[\p{Han}\w]{5,40}`
const UsernamePattern = `[\w]{6,40}`

// PatchJSON ...
type PatchJSON struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func getCount(c *gin.Context) (int64, error) {
	count := c.Query("count")
	if count == "" {
		count = "20"
	}
	intData, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		log.Error(err)
		return 20, err
	}
	return intData, nil
}

func getSkip(c *gin.Context) (int64, error) {
	skip := c.Query("skip")
	if skip == "" {
		skip = "0"
	}
	intData, err := strconv.ParseInt(skip, 10, 64)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return intData, nil
}

func getPatchBody(c *gin.Context) []PatchJSON {
	res := []PatchJSON{}
	err := c.ShouldBindBodyWith(&res, binding.JSON)
	if err != nil {
		log.Error(err)
		return []PatchJSON{}
	}
	// add json validate
	for _, v := range res {
		if v.Op == "" || v.Path == "" || v.Value == "" {
			log.Error("Some patch body format error")
			return []PatchJSON{}
		}
	}
	return res
}

func isEmailVaild(email string) bool {
	reg := regexp.MustCompile(EmailPattern)
	return reg.MatchString(email)
}

func isUsernameValid(username string) bool {
	reg := regexp.MustCompile(UsernamePattern)
	return reg.MatchString(username)
}

func getIntParam(c *gin.Context, name string) (int64, error) {
	idStr := c.Params.ByName(name)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return id, nil
}

func getStringParam(c *gin.Context, name string) (string, error) {
	return c.Params.ByName(name), nil
}

func getCurrentUser(c *gin.Context) models.User {
	return c.MustGet("currentUser").(models.User)
}

func struct2QueryString(obj interface{}) string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	log.Info("tttttt", t)
	log.Info("vvvvvv", v)
	s := "?"
	for i := 0; i < t.NumField(); i++ {
		value := v.Field(i).Interface().(string)

		key := t.Field(i).Name
		jsonKey := ""
		switch key {
		case "UserID":
			jsonKey = "userId"
		case "PlanName":
			jsonKey = "planName"
		case "Frequency":
			jsonKey = "frequency"
		case "SubTime":
			jsonKey = "subTime"
		case "IsActive":
			jsonKey = "isActive"
		case "Key":
			jsonKey = "key"
		case "Secret":
			jsonKey = "secret"
		case "Phrase":
			jsonKey = "phrase"
		case "APIKeyName":
			jsonKey = "apiKeyName"
		case "PlanID":
			jsonKey = "planId"
		case "Exchange":
			jsonKey = "exchange"
		case "Symbol":
			jsonKey = "symbol"
		case "Price":
			jsonKey = "price"
		case "Amount":
			jsonKey = "amount"
		case "Budget":
			jsonKey = "budget"
		case "TxTime":
			jsonKey = "txTime"
		case "MsgType":
			jsonKey = "msgType"
		default:
			continue
		}
		if i != 0 {
			s = s + "&"
		}
		s += jsonKey
		s += "="
		s += value

	}
	log.Info("sssssss", s)
	return s
}

// Struct2Map ...Struct2Map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
