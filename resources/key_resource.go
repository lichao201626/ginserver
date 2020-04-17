package resources

import (
	"ginserver/dao"
	"ginserver/models"
	"ginserver/serializers"
	// "ginserver/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// KeyResource ...
type KeyResource struct {
}

// NewKeyResource ...
func NewKeyResource(e *gin.Engine) {
	u := KeyResource{}
	// Setup Routes
	e.GET("/get/exchanges", u.getKeys)
	e.GET("/get/exchanges/:id", u.getKeysByID)
	e.POST("/post/exchanges", u.postKeys)
	e.PATCH("/patch/exchanges/:id", u.patchKeysByID)
	e.DELETE("/delete/exchanges/:id", u.deleteKeysByID)
}

func (r *KeyResource) postKeys(c *gin.Context) {
	key := models.Key{}
	c.MustBindWith(&key, binding.JSON)
	currentUserID := c.MustGet("currentUserID").(string)
	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	if key.Exchange == "" || key.Key == "" || key.Secret == "" {
		c.JSON(400, serializers.SerializeError(40001, "Must supply exchange, key, secret and phrase"))
		return
	}
	if key.UserID != "" && key.UserID != currentUserID {
		c.JSON(401, serializers.SerializeError(40102, "Can not post key of other user"))
		return
	}
	if key.UserID == "" {
		key.UserID = currentUserID
	}

	// check if the key name exist
	query := bson.M{
		"userId":     currentUserID,
		"apiKeyName": key.APIKeyName,
	}
	total, err := dao.GetKeysTotal(dbSession, query)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Mongo get keys failed"))
		return
	}
	if total > 0 {
		c.JSON(200, serializers.SerializeError(20001, "Key name already exist"))
		return
	}
	// todo verify if this key is useful, call remote method
	pJSON := planJSON{
		UserID:     currentUserID,
		APIKeyName: key.APIKeyName,
		Key:        key.Key,
		Secret:     key.Secret,
		Phrase:     key.Phrase,
	}
	res, validateErr := postValidate(pJSON)
	if validateErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Validate key failed"))
		return
	}
	validateJSON := res.Body.(map[string]interface{})
	if validateJSON["isKeyValid"].(string) != "true" {
		c.JSON(200, serializers.SerializeError(20002, "This key is invalid"))
		return
	}

	encryptErr := key.Encrypt()
	if encryptErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Encrypt the key failed"))
		return
	}

	keyRes, createErr := dao.CreateKey(dbSession, key)
	if createErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo create Keys failed"))
		return
	}
	log.Info("New key created ", keyRes)
	c.JSON(200, serializers.SerializeKey(keyRes))
}

func (r *KeyResource) getKeys(c *gin.Context) {
	currentUserID := c.MustGet("currentUserID").(string)
	userID := c.Query("userId")
	if userID != "" && userID != currentUserID {
		c.JSON(401, serializers.SerializeError(40102, "Can not get keys of other user"))
		return
	}

	skip, _ := getSkip(c)
	count, _ := getCount(c)
	query := bson.M{
		"userId": currentUserID,
	}

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	total, err := dao.GetKeysTotal(dbSession, query)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Mongo get Keys total failed"))
		return
	}
	keys, getErr := dao.GetKeys(dbSession, count, skip, query)
	if getErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo get Keys failed"))
		return
	}
	// decrypt the keys
	for _, v := range keys {
		decryptErr := v.Decrypt()
		if decryptErr != nil {
			c.JSON(500, serializers.SerializeError(50003, "Decrypt the key failed"))
			return
		}
	}
	log.Info("Get keys success ", total)
	c.JSON(200, serializers.SerializeKeys(keys, count, skip, total))
}

func (r *KeyResource) getKeysByID(c *gin.Context) {
	keyID := c.Param("id")
	currentUserID := c.MustGet("currentUserID").(string)

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	key, err := dao.GetKeysByID(dbSession, keyID)
	if err != nil {
		c.JSON(200, serializers.SerializeError(20001, "Not found by the Key id"))
		return
	}
	if key.UserID != currentUserID {
		c.JSON(401, serializers.SerializeError(40102, "Can not get keys of other user"))
		return
	}
	//decrypt the key
	decryptErr := key.Decrypt()
	if decryptErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Decrypt the key failed"))
		return
	}
	log.Info("Get Keys by id success ", keyID)
	c.JSON(200, serializers.SerializeKey(key))
}

func (r *KeyResource) patchKeysByID(c *gin.Context) {
	keyID := c.Param("id")
	currentUserID := c.MustGet("currentUserID").(string)

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	key, err := dao.GetKeysByID(dbSession, keyID)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50001, "Not found by the Key id"))
		return
	}
	if key.UserID != currentUserID {
		c.JSON(401, serializers.SerializeError(40102, "Can not patch keys of other user"))
		return
	}
	//decrypt the key
	decryptErr := key.Decrypt()
	if decryptErr != nil {
		c.JSON(500, serializers.SerializeError(50002, "Decrypt the key failed"))
		return
	}

	body := getPatchBody(c)
	if len(body) == 0 {
		c.JSON(400, serializers.SerializeError(40001, "Empty patch body, or body invalid format"))
		return
	}
	originKey := key.Key
	originKeyName := key.APIKeyName
	originSecret := key.Secret
	originPhrase := key.Phrase
	applyPatchKeyBody(&key, body)

	if key.APIKeyName != originKeyName {
		// check if the key name exist
		query := bson.M{
			"userId":     currentUserID,
			"apiKeyName": key.APIKeyName,
		}
		total, err := dao.GetKeysTotal(dbSession, query)
		if err != nil {
			c.JSON(500, serializers.SerializeError(50001, "Mongo get keys failed"))
			return
		}
		if total > 0 {
			c.JSON(200, serializers.SerializeError(20001, "Key name already exist"))
			return
		}
	}
	if key.Key != originKey || key.Secret != originSecret || key.Phrase != originPhrase {
		pJSON := planJSON{
			UserID:     currentUserID,
			APIKeyName: key.APIKeyName,
			Key:        key.Key,
			Secret:     key.Secret,
			Phrase:     key.Phrase,
		}
		res, validateErr := postValidate(pJSON)
		if validateErr != nil {
			c.JSON(500, serializers.SerializeError(50003, "Validate key failed"))
			return
		}
		validateJSON := res.Body.(map[string]interface{})
		if validateJSON["isKeyValid"].(string) != "true" {
			c.JSON(200, serializers.SerializeError(20002, "This key is invalid"))
			return
		}
	}

	encryptErr := key.Encrypt()
	if encryptErr != nil {
		c.JSON(500, serializers.SerializeError(50004, "Encrypt the key failed"))
		return
	}

	key, err = dao.UpdateKeysByID(dbSession, keyID, key)
	if err != nil {
		c.JSON(500, serializers.SerializeError(50002, "Mongo update Keys by id failed"))
		return
	}
	log.Info("Patch Key by id success ", keyID)
	c.JSON(200, serializers.SerializeKey(key))
}

func applyPatchKeyBody(key *models.Key, body []PatchJSON) {
	for _, v := range body {
		if v.Op != "update" {
			continue
		}
		switch v.Path {
		case "/exchange":
			key.Exchange = v.Value
		case "/key":
			key.Key = v.Value
		case "/secret":
			key.Secret = v.Value
		case "/phrase":
			key.Phrase = v.Value
		case "/apiKeyName":
			key.APIKeyName = v.Value
		default:
			continue
		}
	}
}

func (r *KeyResource) deleteKeysByID(c *gin.Context) {
	keyID := c.Param("id")
	currentUserID := c.MustGet("currentUserID").(string)

	dbSession := c.MustGet("db").(*mgo.Session).Copy()
	defer dbSession.Close()

	key, err := dao.GetKeysByID(dbSession, keyID)
	if err != nil {
		c.JSON(200, serializers.SerializeError(20001, "Not found by the Key id"))
		return
	}
	if key.UserID != currentUserID {
		c.JSON(401, serializers.SerializeError(40102, "Can not delete keys of other user"))
		return
	}

	//
	planQuery := bson.M{
		"userId":   key.UserID,
		"apiKeyId": key.APIKeyID.Hex(),
	}
	total, totalErr := dao.GetPlansTotal(dbSession, planQuery)
	if totalErr != nil {
		c.JSON(500, serializers.SerializeError(50001, "Get plans by the Key id failed"))
		return
	}
	if total > 0 {
		c.JSON(200, serializers.SerializeError(20001, "Can not remove the key using by a plan"))
		return
	}

	query := bson.M{
		"userId":     key.UserID,
		"apiKeyName": key.APIKeyName,
		"key":        key.Key,
		"secret":     key.Secret,
		"phrase":     key.Phrase,
	}

	// todo check if the keys has plan that is active
	err = dao.RemoveKeys(dbSession, query)
	if err != nil {
		c.JSON(200, serializers.SerializeError(20001, "Remove by the Key id"))
		return
	}
	log.Info("Remove Keys by id success ", keyID)
	c.JSON(200, serializers.SerializeResponse(0, nil, "Success"))
}
