package utils

import (
	"github.com/ccwings/log"
	"github.com/joegrasse/mail"

)

// call example:
// utils.SendMail([]string{"liuyuhang171@163.com","michael@ccwings.cn"},"test1","test hello")
func SendMail(sendTo []string,title, content string){
	log.Debug("In Send Mail")
	htmlBody :=
		`<html>
			<body>
				<p><b>This email from HoldingCloud.</b>.</p>
				<p>`+content+`</p>
				<p>Please NOT Reply.</p>
			</body>
		</html>`
	email := mail.New()

	email.Username = "liuyuhang171@163.com"
	email.Password = "u125478521"
	email.SetFrom("HoldingCloud2.0 <liuyuhang171@163.com>")
	email.SetSubject(title)

	for _,v := range sendTo{
		email.AddTo(v)
	}

	email.AddAlternative("text/html", htmlBody)
	err := email.Send("smtp.163.com:25")
	if err != nil {
		log.Debug(err.Error())
	} else {
		log.Debug("Email Send")
	}
}
