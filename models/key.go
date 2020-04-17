package models

import (
	"ginserver/util"
	"gopkg.in/mgo.v2/bson"
)

const (
	// CollectionKey is the mongo collection name of key
	CollectionKey = "keys"
)

// Key is a record of a person's exchange API key
type Key struct {
	APIKeyID   bson.ObjectId `json:"apiKeyId,omitempty" bson:"_id,omitempty"`
	UserID     string        `json:"userId" bson:"userId"`
	Exchange   string        `json:"exchange" bson:"exchange"`
	Key        string        `json:"key" bson:"key"`
	Secret     string        `json:"secret" bson:"secret"`
	Phrase     string        `json:"phrase" bson:"phrase"`
	APIKeyName string        `json:"apiKeyName" bson:"apiKeyName"`
	CreateTime string        `json:"createTime" bson:"createTime"`
	UpdateTime string        `json:"updateTime" bson:"updateTime"`
}

// Encrypt the key
func (key *Key) Encrypt() error {
	encryptedKey, encryptKeyErr := util.AES128CBCEncrypt(key.Key)
	if encryptKeyErr != nil {
		return encryptKeyErr
	}
	key.Key = encryptedKey

	encryptedSecret, encryptSecretErr := util.AES128CBCEncrypt(key.Secret)
	if encryptSecretErr != nil {
		return encryptSecretErr
	}
	key.Secret = encryptedSecret

	if key.Phrase != "" {
		encryptedPhrase, encryptPhraseErr := util.AES128CBCEncrypt(key.Phrase)
		if encryptPhraseErr != nil {
			return encryptPhraseErr
		}
		key.Phrase = encryptedPhrase
	}
	return nil
}

// Decrypt the key
func (key *Key) Decrypt() error {
	decryptedKey, decryptKeyErr := util.AES128CBCDecrypt(key.Key)
	if decryptKeyErr != nil {
		return decryptKeyErr
	}
	key.Key = decryptedKey

	decryptedSecret, decryptSecretErr := util.AES128CBCDecrypt(key.Secret)
	if decryptSecretErr != nil {
		return decryptSecretErr
	}
	key.Secret = decryptedSecret

	if key.Phrase != "" {
		decryptedPhrase, decryptPhraseErr := util.AES128CBCDecrypt(key.Phrase)
		if decryptPhraseErr != nil {
			return decryptPhraseErr
		}
		key.Phrase = decryptedPhrase
	}
	return nil
}
