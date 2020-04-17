package resources

import (
	"ginserver/dao"
	"ginserver/models"
	"ginserver/serializers"
	"ginserver/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserResource ...
type UserResource struct {
}

// NewUserResource ...
func NewUserResource(e *gin.Engine) {
	u := UserResource{}
	// Setup Routes
	e.GET("/get/users", u.getUsers)
	e.GET("/get/users/:id", u.getUsersByID)
	e.POST("/post/users", u.postUsers)
	e.PATCH("/patch/users/:id", u.patchUsersByID)
}

func (r *UserResource) postUsers(c *gin.Context) {
	user := postJSON{}
	c.MustBindWith(&user, binding.JSON)

	verifyEmail := c.Query("verifyEmail")
	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()
	if user.Email == "" && user.Username == "" {
		c.JSON(400, serializers.SerializeError(40001, "Must supply email or username"))
		return
	}
	if user.Email != "" && !isEmailVaild(user.Email) {
		c.JSON(200, serializers.SerializeError(20001, "Email format is invalid"))
		return
	}
	if user.Username != "" && !isUsernameValid(user.Username) {
		c.JSON(200, serializers.SerializeError(20002, "Username format is invalid"))
		return
	}
	// check if user exist
	query := bson.M{}
	if user.Email != "" {
		query["email"] = user.Email
	}
	if user.Username != "" {
		query["username"] = user.Username
	}
	users, err := dao.GetUsers(dbSession, 20, 0, query)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Mongo get users failed"))
		return
	}
	if len(users) > 0 {
		c.JSON(201, serializers.SerializeUser(users[0]))
		if users[0].Password == "" {
			emailPinCode(dbSession, &users[0], "SignUp")
		}
		return
	}

	if user.Email == "" {
		c.JSON(200, serializers.SerializeError(20003, "Username not registered, please supply a email to register"))
		return
	}
	userM := models.User{
		Username: user.Username,
		Email:    user.Email,
	}
	userRes, createErr := dao.CreateUser(dbSession, userM)
	if createErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo create users failed"))
		return
	}
	log.Info("New user created ", userRes)
	if verifyEmail != "false" {
		emailPinCode(dbSession, &userRes, "SignUp")
	}
	c.JSON(200, serializers.SerializeUser(userRes))
}

func (r *UserResource) getUsers(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	if currentUser.IsAdmin != "true" {
		c.JSON(401, serializers.SerializeError(40102, "Can not get info of other user"))
		return
	}

	skip, _ := getSkip(c)
	count, _ := getCount(c)
	createBefore := c.Query("createBefore")
	createAfter := c.Query("createAfter")
	username := c.Query("username")
	email := c.Query("email")

	query := bson.M{}
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
	if username != "" {
		query["username"] = username
	}
	if email != "" {
		query["email"] = email
	}

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	total, err := dao.GetUsersTotal(dbSession, query)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Mongo get users total failed"))
		return
	}
	users, getErr := dao.GetUsers(dbSession, count, skip, query)
	if getErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo get users failed"))
		return
	}
	log.Info("Get users success ", total)
	c.JSON(200, serializers.SerializeUsers(users, count, skip, total))
}

func (r *UserResource) getUsersByID(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.MustGet("currentUserID").(string)
	currentUser := c.MustGet("currentUser").(models.User)

	if userID != currentUserID && currentUser.IsAdmin != "true" {
		c.JSON(401, serializers.SerializeError(40102, "Can not get info of other user"))
		return
	}

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	user, err := dao.GetUsersByID(dbSession, userID)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50003, "Not found by the user id"))
		return
	}
	log.Info("Get users by id success ", userID)
	c.JSON(200, serializers.SerializeUser(user))
}

func (r *UserResource) patchUsersByID(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.MustGet("currentUserID").(string)
	onlyChPwd := ""
	onlyChPwdInterface, _ := c.Get("onlyChPwd")
	if onlyChPwdInterface != nil {
		onlyChPwd = "true"
	}

	if userID != currentUserID && onlyChPwd == "" {
		c.JSON(401, serializers.SerializeError(40101, "Can not change info of other user"))
		return
	}

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	user, err := dao.GetUsersByID(dbSession, userID)
	if err != nil {
		c.JSON(200, serializers.SerializeError(20001, "Not found by the user id"))
		return
	}

	body := getPatchBody(c)
	if len(body) == 0 {
		c.JSON(400, serializers.SerializeError(40001, "Empty patch body, or body invalid format"))
		return
	}

	userName := user.Username
	applyPatchBody(&user, body, onlyChPwd)

	if user.Username != "" && !isUsernameValid(user.Username) {
		c.JSON(200, serializers.SerializeError(20002, "Username format not valid"))
		return
	}
	if len(user.Nickname) > 100 {
		c.JSON(200, serializers.SerializeError(20003, "Nickname is too long"))
		return
	}
	if userName != user.Username {
		queryUsername := bson.M{
			"username": user.Username,
		}
		totalUsername, getErr := dao.GetUsersTotal(dbSession, queryUsername)
		if getErr != nil {
			c.JSON(500, serializers.SerializeError(50001, "Mongo get users total failed"))
			return
		}
		if totalUsername != 0 {
			c.JSON(200, serializers.SerializeError(20004, "Username already used, try a new one"))
			return
		}
	}

	user, err = dao.UpdateUsersByID(dbSession, userID, user)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo update users by id failed"))
		return
	}
	log.Info("Patch user by id success ", userID)
	c.JSON(200, serializers.SerializeUser(user))
}

func applyPatchBody(user *models.User, body []PatchJSON, onlyChPwd string) {
	for _, v := range body {
		if v.Op != "update" {
			continue
		}
		if onlyChPwd == "true" {
			switch v.Path {
			case "/password":
				user.Password = util.Md5String(user.CreateTime + v.Value)
			default:
				continue
			}
		} else {
			switch v.Path {
			case "/username":
				user.Username = v.Value
			case "/nickname":
				user.Nickname = v.Value
			case "/password":
				user.Password = util.Md5String(user.CreateTime + v.Value)
			default:
				continue
			}
		}
	}
}
