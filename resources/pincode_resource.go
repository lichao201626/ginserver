package resources

import (
	"ginserver/dao"
	"ginserver/models"
	"ginserver/serializers"
	"ginserver/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// PincodeResource ...
type PincodeResource struct {
}

// postJSON ...
type postJSON struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Pincode  string `json:"pincode"`
	Password string `json:"password"`
}

// NewPincodeResource ...
func NewPincodeResource(e *gin.Engine) {
	u := PincodeResource{}
	e.POST("/post/pincode", u.postPincode)
	e.POST("/post/verify", u.verifyPincode)
}

func (r *PincodeResource) postPincode(c *gin.Context) {
	user := postJSON{}
	c.MustBindWith(&user, binding.JSON)

	if user.Username == "" && user.Email == "" {
		c.JSON(400, serializers.SerializeError(40001, "Must supply email or username"))
		return
	}

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	query := bson.M{}
	if user.Username != "" {
		query["username"] = user.Username
	}
	if user.Email != "" {
		query["email"] = user.Email
	}
	users, err := dao.GetUsers(dbSession, 20, 0, query)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Mongo get users failed"))
		return
	}
	if len(users) != 1 {
		c.JSON(200, serializers.SerializeError(20001, "User not rejistered, please register first"))
		return
	}

	total, emailErr := emailPinCode(dbSession, &users[0], "ForgetPsw")
	if emailErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo failed when email pincode"))
		return
	}
	if total > 0 {
		c.JSON(201, serializers.SerializeError(20101, "Do not try in less 15 minute"))
		return
	}
	c.JSON(200, serializers.SerializeResponse(0, nil, "Success"))
}

func emailPinCode(dbSession *mgo.Session, user *models.User, reason string) (int, error) {
	expire := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	expireQuery := bson.M{
		"expire": bson.M{
			"$lte": expire,
		},
	}
	err := dao.RemovePincodes(dbSession, expireQuery)
	if err != nil {
		return 0, err
	}
	total, getErr := dao.GetPincodesTotal(dbSession, bson.M{
		"userId": user.UserID.Hex(),
	})
	if getErr != nil {
		return 0, getErr
	}
	if total > 0 {
		return total, nil
	}
	pinCode := util.RandString(6, 0)
	Pincode := models.Pincode{}
	Pincode.PinCode = util.Md5String(user.CreateTime + pinCode)
	Pincode.UserID = user.UserID.Hex()
	_, createErr := dao.CreatePincode(dbSession, Pincode)
	if createErr != nil {
		return 0, createErr
	}
	go util.MailToUser(user.Email, reason, pinCode)
	return 0, nil
}

func (r *PincodeResource) verifyPincode(c *gin.Context) {
	verifyJSON := postJSON{}
	c.MustBindWith(&verifyJSON, binding.JSON)

	if verifyJSON.Username == "" && verifyJSON.Email == "" {
		c.JSON(400, serializers.SerializeError(40001, "Must supply email or username"))
		return
	}

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	query := bson.M{}
	if verifyJSON.Username != "" {
		query["username"] = verifyJSON.Username
	}
	if verifyJSON.Email != "" {
		query["email"] = verifyJSON.Email
	}
	users, err := dao.GetUsers(dbSession, 20, 0, query)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Mongo get users failed"))
		return
	}
	if len(users) != 1 {
		c.JSON(200, serializers.SerializeError(20001, "User not rejistered, please register first"))
		return
	}

	expire := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	PincodeQuery := bson.M{
		"userId": users[0].UserID.Hex(),
		"expire": bson.M{
			"$gte": expire,
		},
		"pincode": util.Md5String(users[0].CreateTime + verifyJSON.Pincode),
	}
	Pincodes, pincodeErr := dao.GetPincodes(dbSession, 20, 0, PincodeQuery)
	if pincodeErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo get pincodes failed"))
		return
	}
	if len(Pincodes) == 0 {
		c.JSON(200, serializers.SerializeError(20002, "Wrong pincode or expired, try new pincode"))
		return
	}
	if Pincodes[0].Tried >= 3 {
		c.JSON(200, serializers.SerializeError(20003, "Wrong pincode 3 times, cannot retry within 15 minutes"))
		return
	}

	pincodeID := Pincodes[0].PincodeID.Hex()
	Pincodes[0].Tried++
	dao.UpdatePincodesByID(dbSession, pincodeID, Pincodes[0])

	authToken, tokenErr := newToken(dbSession, &users[0], "ForgetPsw")
	if tokenErr != nil {
		c.JSON(500, serializers.SerializeError(50003, "Mongo new token failed"))
		return
	}
	token := models.Token{
		UserID:    users[0].UserID.Hex(),
		AuthToken: authToken,
	}
	c.JSON(200, serializers.SerializeToken(token))
}
