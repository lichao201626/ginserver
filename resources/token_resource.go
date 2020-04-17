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
)

// TokenResource ...
type TokenResource struct {
}

// NewTokenResource ...
func NewTokenResource(e *gin.Engine) {
	u := TokenResource{}
	e.POST("/post/login", u.postLogin)
	e.POST("/post/logout", u.postLogout)
}

func (r *TokenResource) postLogin(c *gin.Context) {
	verifyJSON := postJSON{}
	c.MustBindWith(&verifyJSON, binding.JSON)

	if verifyJSON.Username == "" && verifyJSON.Email == "" {
		c.JSON(400, serializers.SerializeError(40001, "Must supply email or username"))
		return
	}
	if verifyJSON.Password == "" {
		c.JSON(400, serializers.SerializeError(40002, "Must supply a passowrd"))
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

	queryByUserID := bson.M{
		"userId": users[0].UserID.Hex(),
	}
	removeErr := dao.RemoveTokens(dbSession, queryByUserID)
	if removeErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo remove tokens failed"))
		return
	}

	pass := util.Md5String(users[0].CreateTime + verifyJSON.Password)
	if pass != users[0].Password {
		c.JSON(200, serializers.SerializeError(20002, "Wrong password"))
		return
	}

	authToken, tokenErr := newToken(dbSession, &users[0], "login")
	if tokenErr != nil {
		c.JSON(500, serializers.SerializeError(50003, "Mongo create tokens failed"))
		return
	}
	token := models.Token{
		UserID:    users[0].UserID.Hex(),
		AuthToken: authToken,
	}
	c.JSON(200, serializers.SerializeToken(token))
}

func (r *TokenResource) postLogout(c *gin.Context) {
	token := models.Token{}
	c.Bind(&token)

	if token.UserID == "" {
		c.JSON(400, serializers.SerializeError(40001, "Must supply your userId"))
		return
	}
	currentUserID := c.MustGet("currentUserID").(string)
	if token.UserID != currentUserID {
		c.JSON(401, serializers.SerializeError(40102, "Can not logout other user"))
		return
	}

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	query := bson.M{
		"userId": token.UserID,
	}
	err := dao.RemoveTokens(dbSession, query)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Mongo remove tokens failed"))
		return
	}
	c.JSON(200, serializers.SerializeResponse(0, nil, "Success"))
}

func newToken(dbSession *mgo.Session, user *models.User, tokenType string) (string, error) {
	temp := util.RandString(10, 3)
	tokenString := util.Md5String(temp + user.UserID.Hex() + user.CreateTime)

	authToken := util.Md5String("token" + tokenString)
	token := models.Token{}
	token.AuthToken = authToken
	token.UserID = user.UserID.Hex()
	token.Type = tokenType

	_, err := dao.CreateToken(dbSession, token)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
