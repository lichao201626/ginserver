package db

import (
	"gopkg.in/mgo.v2"
)

// Session Client creates new client based on the ClientConfig provided.
func Session() *mgo.Session {
	// 用户密码 mongodb://foo:bar@localhost:27017
	s, err := mgo.Dial("mongodb://172.17.0.1:27017")

	// Check if connection error, is mongo running
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}
