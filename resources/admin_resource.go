package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"ginserver/dao"
	"ginserver/models"
	"ginserver/serializers"
	"ginserver/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AdminResource ...
type AdminResource struct {
}

// NewAdminResource ...
func NewAdminResource(e *gin.Engine) {
	u := AdminResource{}
	e.POST("/admin/login", u.postLogin)
	e.POST("/admin/logout", u.postLogout)
}

func (r *AdminResource) postLogin(c *gin.Context) {
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
	temp := util.RandString(10, 3)
	authToken := util.Md5String(temp + verifyJSON.Username + verifyJSON.Password)

	token := models.Token{
		AuthToken: authToken,
	}
	c.JSON(200, serializers.SerializeToken(token))
}

func (r *AdminResource) postLogout(c *gin.Context) {
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
