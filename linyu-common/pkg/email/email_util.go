package email

import (
	"bytes"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/logger"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"html/template"
)

var EmailDir string

func Init(emailDir string) {
	EmailDir = emailDir
}

// SendEmailCode 发送验证码邮件
func SendEmailCode(to string, code string) {
	go func() {
		err := SendEmail([]string{to}, map[string]string{"code": code}, "林语验证码: "+code, "code.html")
		if err != nil {
			logger.Log.Error("[SendEmailCode] 邮件发送失败:", zap.Error(err))
		}
	}()
}

// SendEmail 发送邮件
func SendEmail(to []string, data interface{}, subject string, tmplName string) error {
	tmplPath := EmailDir + "/" + tmplName
	// 解析模板文件
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}
	// 渲染模板内容
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}
	// 构建邮件
	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress(config.C.Email.FormAddr, config.C.Email.Username))
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body.String())
	// 发送
	d := gomail.NewDialer(config.C.Email.Host, config.C.Email.Port, config.C.Email.FormAddr, config.C.Email.Password)
	if err := d.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
