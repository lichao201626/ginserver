package util

import (
	"ginserver/util"
	"testing"
)

func Test_MailToUser(t *testing.T) {
	/* 	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} */
	isSend := util.MailToUser("824683639@qq.com", "ForgetPsw", "123452")
	if isSend != true {
		//	t.Error("af")
	}
}
