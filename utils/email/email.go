package email

import (
	"bytes"
	"html/template"

	"github.com/playground-pro-project/playground-pro-api/app/config"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/pquerna/otp/totp"
	"gopkg.in/gomail.v2"
)

var log = middlewares.Log()

func SendOTP(name, email string) {
	var body bytes.Buffer
	t := template.New("otp.html")
	t, err := t.Parse(`
		<html>
			<body>
				<h3>Kode OTP Anda</h3>
				<p>Hallo {{.Name}}, Apa kabar? Semoga harimu menyenangkan!</p>
				<p>Kode OTP: {{.OTPCode}}</p>
			</body>
		</html>
	`)
	if err != nil {
		log.Error("error parsing template email")
		return
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      config.GOMAIL_HOST,
		AccountName: config.GOMAIL_EMAIL,
		SecretSize:  7,
	})

	if err != nil {
		panic(err)
	}
	err = t.Execute(&body, struct {
		Name    string
		OTPCode string
	}{
		Name:    name,
		OTPCode: key.Secret(),
	})
	if err != nil {
		log.Error("error rendering template email")
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.GOMAIL_EMAIL)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Kode OTP Anda")
	m.SetBody("text/html", body.String())
	d := gomail.NewDialer(config.GOMAIL_HOST, config.GOMAIL_PORT, config.GOMAIL_EMAIL, config.GOMAIL_PASSWORD)
	if err := d.DialAndSend(m); err != nil {
		log.Sugar().Error("Gagal mengirim email: ", err.Error())
		return
	}

	log.Info("Email terkirim.")
}
