package resources

import (
	// "encoding/json"
	"fmt"
	"ginserver/dao"
	"ginserver/models"
	"ginserver/serializers"
	"ginserver/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

// OrderResource ...
type OrderResource struct {
}

// NewOrderResource ...
func NewOrderResource(e *gin.Engine) {
	u := OrderResource{}
	// Setup Routes
	e.GET("/get/orders", u.getOrders)
	e.GET("/get/orders/:id", u.getOrdersByID)
	e.POST("/post/orders", u.postOrders)
	e.POST("/post/email", u.postEmail)
	e.POST("/post/validate", u.postValidate)
	// e.PATCH("/patch/Orders/:id", u.patchOrdersByID)
}

//OrderJSON ..
type OrderJSON struct {
	PlanID   string `json:"planId"`
	UserID   string `json:"userId"`
	Budget   string `json:"budget"`
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
	Price    string `json:"price"`
	Amount   string `json:"amount"`
	TxTime   string `json:"txTime"`
}

type emailJSON struct {
	UserID  string `json:"userId"`
	MsgType string `json:"msgType"`
}

func (r *OrderResource) postOrders(c *gin.Context) {
	order := OrderJSON{}
	c.MustBindWith(&order, binding.JSON)
	isReqPass := reqAuthentication(c, order)
	if !isReqPass {
		return
	}
	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	currentUser, err := dao.GetUsersByID(dbSession, order.UserID)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Get the user of the order failed"))
		return
	}

	// todo get the userID by the plan id
	orderM := models.Order{
		PlanID:    order.PlanID,
		UserID:    order.UserID,
		Exchange:  order.Exchange,
		Price:     order.Price,
		Unit:      order.Symbol,
		Amount:    order.Amount,
		Fee:       order.Budget,
		OrderTime: order.TxTime,
	}
	orderRes, createErr := dao.CreateOrder(dbSession, orderM)
	if createErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo create Orders failed"))
		return
	}
	log.Info("New Order created ", orderRes)
	c.JSON(200, serializers.SerializeOrder(orderRes))
	go util.MailToUser(currentUser.Email, "OrderSuccessed", "")
}

func (r *OrderResource) getOrders(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	/* 	currentOrder := c.MustGet("currentOrder").(models.Order)
	   	if currentOrder.IsAdmin != "true" {
	   		c.JSON(401, serializers.SerializeError(40102, "Can not get info of other Order"))
	   		return
	   	} */

	skip, _ := getSkip(c)
	count, _ := getCount(c)
	createBefore := c.Query("createBefore")
	createAfter := c.Query("createAfter")
	// ordername := c.Query("ordername")
	// email := c.Query("email")
	// log.Info("Ordername:", ordername)

	query := bson.M{}
	query["userId"] = currentUser.UserID.Hex()
	if createAfter != "" && createBefore == "" {
		query["orderTime"] = bson.M{
			"$gte": createAfter,
		}
	}
	if createBefore != "" && createAfter == "" {
		query["orderTime"] = bson.M{
			"$lte": createBefore,
		}
	}
	if createBefore != "" && createAfter != "" {
		query["createTime"] = bson.M{
			"$gte": createAfter,
			"$lte": createBefore,
		}
	}

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	total, err := dao.GetOrdersTotal(dbSession, query)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Mongo get Orders total failed"))
		return
	}
	orders, getErr := dao.GetOrders(dbSession, count, skip, query)
	if getErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo get Orders failed"))
		return
	}
	log.Info("Get orders success ", total)
	c.JSON(200, serializers.SerializeOrders(orders, count, skip, total))
}

func (r *OrderResource) getOrdersByID(c *gin.Context) {
	orderID := c.Param("id")
	currentUser := c.MustGet("currentUser").(models.User)

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	order, err := dao.GetOrdersByID(dbSession, orderID)
	if err != nil {
		c.JSON(200, serializers.SerializeError(20001, "Not found by the order id"))
		return
	}
	if order.UserID != currentUser.UserID.Hex() {
		c.JSON(200, serializers.SerializeError(20001, "Can not get order of other user"))
		return
	}
	log.Info("Get orders by id success ", orderID)
	c.JSON(200, serializers.SerializeOrder(order))
}

func (r *OrderResource) postEmail(c *gin.Context) {
	eJSON := emailJSON{}
	c.MustBindWith(&eJSON, binding.JSON)
	// currentUser := c.MustGet("currentUser").(models.User)
	reqAuthentication(c, eJSON)

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	currentUser, err := dao.GetUsersByID(dbSession, eJSON.UserID)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Get the user of the order failed"))
		return
	}
	c.JSON(200, serializers.SerializeResponse(0, nil, "Success"))
	go util.MailToUser(currentUser.Email, "OrderSuccessed", "")
}

func encryptReq(uri string, method string, reqJSON interface{}) (string, string) {
	apiKey := "quanta_lab_aip_"
	secret := "test"
	timestamp := time.Now().Unix()
	timestring := strconv.Itoa(int(timestamp))

	key := apiKey + timestring
	log.Info("key::", key)
	fmt.Println(key)
	prehash := timestring + method + uri + struct2QueryString(reqJSON)
	log.Info("prehash::", prehash)
	fmt.Println(prehash)
	sign := util.ComputeHmac256(prehash, secret)
	log.Info("sign::", sign)

	return key, sign
}

func reqAuthentication(c *gin.Context, reqJSON interface{}) bool {
	log.Info(c.Request.Header)
	key := c.Request.Header.Get("BTAI-ACCESS-KEY")
	if key == "" {
		c.JSON(401, serializers.SerializeError(40102, "Access key not exist"))
		return false
	}
	sign := c.Request.Header.Get("BTAI-ACCESS-SIGN")
	if sign == "" {
		c.JSON(401, serializers.SerializeError(40103, "Access sign not exist"))
		return false
	}

	apiKey := "quanta_lab_aip_"
	if !strings.HasPrefix(key, apiKey) {
		c.JSON(401, serializers.SerializeError(40104, "Access key format invalid"))
		return false
	}
	timestring := key[15:]
	fmt.Println(timestring)
	/* 	timestamp, atoiErr := strconv.Atoi(timestring)
	   	if atoiErr != nil {
	   		c.JSON(401, "Access key format not valid")
	   		return
	   	}
	   	// within 2 min
	   	//timestamp := time.Now().Unix()

	   	beforTime := time.Now().Add(-time.Minute * 1).Unix()
	   	afterTime := time.Now().Add(time.Minute * 1).Unix()
	   	if int64(timestamp) > afterTime || int64(timestamp) < beforTime {
	   		c.JSON(401, "Time stamp out of range")
	   		return
	   	}
	*/
	secret := "test"

	method := strings.ToLower(c.Request.Method)
	uri := c.Request.RequestURI

	log.Info("reqjson:", reqJSON)
	prehash := timestring + method + uri + struct2QueryString(reqJSON)
	log.Info("prehash:", prehash)
	fmt.Println(prehash)
	encryptSign := util.ComputeHmac256(prehash, secret)
	log.Info("encryptSign", encryptSign)

	if sign != encryptSign {
		c.JSON(401, serializers.SerializeError(40105, "Access sign format invalid"))
		return false
	}

	return true
}

func (r *OrderResource) postValidate(c *gin.Context) {
	// order := OrderJSON{}
	order := planJSON{}
	c.MustBindWith(&order, binding.JSON)
	ra := reqAuthentication(c, order)
	if !ra {
		return
	}
	vs := ValidateRes{
		IsKeyValid:     "true",
		IsPlanValid:    "true",
		IsBudgetEnough: "true",
	}

	c.JSON(200, serializers.SerializeResponse(0, vs, "Success"))
	// go util.MailToUser(currentUser.Email, "OrderSuccessed", "")
}
