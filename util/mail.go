package util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/mail"
	"net/smtp"
)

// MailToUser ...
/*

switch contentType {
case "SignUp":
	发送注册验证码
case "ForgetPsw":
	发送忘记密码验证码
case "OrderSuccessed":
	发送下单成功通知
case "OrderFailed":
	发送由于资金不足下单失败通知
case "FundInsufficient":
	发送资金不足提前通知
}

*/
func MailToUser(address string, contentType string, content string) bool {
	mailer := "yawei.hong@invault.io"
	passPhrase := "y93U6aGrARvKrjfR"
	servername := "smtphm.qiye.163.com:994"
	from := mail.Address{Name: "币投宝", Address: mailer}
	to := mail.Address{Name: "", Address: address}

	// 根据contentType调整邮件内容
	var body bytes.Buffer

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Content-Type"] = "text/html;charset=UTF-8"

	switch contentType {
	case "SignUp":
		headers["Subject"] = "【币投宝邮箱验证】"
		for k, v := range headers {
			body.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		}
		// 发送注册验证码
		temp, _ := template.ParseFiles("templates/signUpMail.html")
		temp.Execute(&body, struct {
			Address string
			Pin     string
		}{
			Address: address,
			Pin:     content,
		})
	case "ForgetPsw":
		headers["Subject"] = "【币投宝邮箱验证】"
		for k, v := range headers {
			body.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		}
		// 发送忘记密码验证码
		temp, _ := template.ParseFiles("templates/forgetPswMail.html")
		temp.Execute(&body, struct {
			Address string
			Pin     string
		}{
			Address: address,
			Pin:     content,
		})
	case "OrderSuccessed":
		headers["Subject"] = "【您的定投计划执行成功】"
		for k, v := range headers {
			body.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		}
		// 发送下单成功通知
		temp, _ := template.ParseFiles("templates/orderSuccessedMail.html")
		temp.Execute(&body, struct{}{})
	case "OrderFailed":
		headers["Subject"] = "【您的定投执行失败】"
		for k, v := range headers {
			body.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		}
		// 发送由于资金不足下单失败通知
		temp, _ := template.ParseFiles("templates/orderFailedMail.html")
		temp.Execute(&body, struct{}{})
	case "FundInsufficient":
		headers["Subject"] = "【您的交易所账号余额不足】"
		for k, v := range headers {
			body.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		}
		// 发送资金不足提前通知
		temp, _ := template.ParseFiles("templates/fundInsufficientMail.html")
		temp.Execute(&body, struct{}{})
	}

	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", mailer, passPhrase, host)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	isSent := true
	// recover when unexpect panic
	defer func() {
		recover()
		isSent = false
	}()

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}
	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}
	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write(body.Bytes())
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
	return isSent
}
