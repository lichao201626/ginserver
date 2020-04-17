package util

import (
	"ginserver/util"
	// "github.com/stretchr/testify/assert"
	//"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

type ReqJSON struct {
	PlanId string `json:"planId"`
	// Exchange string `json:"exchange"`
}

func Test_Hmac256(t *testing.T) {

	/* 	r := gin.Default()
	   	resources.NewUserResource(r)
	   	w := httptest.NewRecorder()
	   	req, _ := http.NewRequest("GET", "/get/users", nil)
	   	r.ServeHTTP(w, req)

	   	assert.Equal(t, 500, w.Code)
	   	t.Error("case fail", w.Body.String()) */
	// assert.Equal(t, "pong", w.Body.String())
	timestamp := time.Now().Unix()
	// timestamp := 1546844474
	// t.Error(strconv.Itoa(timestamp))
	//1546844474
	method := "post"
	uri := "/post/orders"
	secret := "test"
	body := ReqJSON{
		PlanId: "gfds",
		// Exchange: "gfds",
	}
	//kC37Z+gA3/Z0WispVSDraZNaWpPJi2QPK0rskOBkLKM=

	t.Error(body)
	//m := Struct2Map(body)
	//s, err := json.Marshal(m)
	//body1 := ReqJSON{}
	//t.Error(string(s), err)

	//json.Unmarshal(s, &body1)
	//t.Error(body1)
	//string(timestamp)
	// intData, err := strconv.ParseInt(timestamp, 10, 64)
	prehash := strconv.Itoa(int(timestamp)) + method + uri + Struct2QueryString(body)
	prehash = "1546868386post/post/orders?planId=gfds"

	// prehash :=strconv.Itoa(timestamp) +method + uri + "?" + string(s)
	// prehash = "1546844474post/post/abc?id=asdf&name=wdc"
	t.Error(prehash)
	x := util.ComputeHmac256(prehash, secret)
	t.Error(x)
}
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func Struct2QueryString(obj interface{}) string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	s := "?"
	for i := 0; i < t.NumField(); i++ {
		if i != 0 {
			s += "&"
		}
		s += strings.ToLower(t.Field(i).Name)
		s += "="
		s += v.Field(i).Interface().(string)
	}
	return s
}
