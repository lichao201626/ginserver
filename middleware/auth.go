package middleware

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"ginserver/dao"
	"ginserver/util"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Error response format
type Error struct {
	Code    int         `json:"code"`
	Body    interface{} `json:"body"`
	Message string      `json:"message"`
}

func newAuthError() Error {
	return Error{
		Code:    40101,
		Message: "Unauthorized, please login",
	}
}

// Auth ......
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		r, _ := regexp.Compile("^Token (.+)$")
		match := r.FindStringSubmatch(authHeader)

		noAuth := []string{
			"/admin/login",
			"/admin/logout",
			"/post/users",
			"/post/login",
			"/post/pincode",
			"/post/verify",
			"/post/priceData",
			"/post/orders",
			"/post/validate",
			"/ws",
		}

		for _, v := range noAuth {
			fmt.Println("v", v, c.Request.RequestURI)
			if v == c.Request.RequestURI {
				c.Next()
				return
			}
		}

		authError := newAuthError()
		if len(match) == 0 {
			c.AbortWithStatusJSON(401, authError)
			return
		}
		tokenString := match[1]

		if len(tokenString) == 0 {
			c.AbortWithStatusJSON(401, authError)
			return
		}

		dbSession := c.MustGet("db").(*mgo.Session).Copy()
		defer dbSession.Close()

		expire := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		expiredQuery := bson.M{
			"expire": bson.M{
				"$lte": expire,
			},
		}
		removeErr := dao.RemoveTokens(dbSession, expiredQuery)
		if removeErr != nil {
			c.AbortWithStatusJSON(401, Error{
				Code:    40111,
				Message: "Mongo remove token failed",
			})
			return
		}

		query := bson.M{
			"authToken": util.Md5String("token" + tokenString),
		}
		tokens, getErr := dao.GetTokens(dbSession, 20, 0, query)
		if getErr != nil {
			c.AbortWithStatusJSON(401, Error{
				Code:    40112,
				Message: "Mongo get tokens failed",
			})
			return
		}
		if len(tokens) != 1 {
			c.AbortWithStatusJSON(401, authError)
			return
		}

		userID := tokens[0].UserID
		user, err := dao.GetUsersByID(dbSession, userID)
		if err != nil {
			c.AbortWithStatusJSON(401, Error{
				Code:    40113,
				Message: "Mongo get users by token userId failed",
			})
			return
		}

		if tokens[0].Type == "ForgetPsw" {
			if !strings.HasPrefix(c.Request.RequestURI, "/patch/users/") {
				c.AbortWithStatusJSON(401, authError)
				return
			}
			c.Set("onlyChPwd", "true")
		}

		c.Set("currentUser", user)
		c.Set("currentUserID", user.UserID.Hex())
		c.Next()
	}
}
