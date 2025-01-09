package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/sabt-dev/0-Project/internal/config"
	"github.com/thanhpk/randstr"
	"gopkg.in/gomail.v2"
)

func ValidateEmailFormat(email string) error {
	// Check email format
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

// VerifyEmailExistence checks if the email format is valid and if the domain has MX records
func VerifyEmailExistence(email string) (bool, error) {
	// Check email format
	err := ValidateEmailFormat(email)
	if err != nil {
		return false, err
	}
	// Check domain MX records
	domain := email[strings.LastIndex(email, "@")+1:]
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return false, errors.New("domain does not have MX records")
	}

	return true, nil
}

// SendVerificationCode sends the verification code to the user's email
func SendVerificationCode(email, code, firstName string) error {
	//TODO: Implement the email sending logic to .env variables
	from := config.AppConfig.FromEmail          //"your-email@example.com"
	smtpPass := config.AppConfig.SMTPPass  //"your-email-password"
	smtpUser := config.AppConfig.SMTPUser  //"your-email-username"
	to := email
	smtpHost := config.AppConfig.SMTPHost     //"smtp.gmail.com"
	smtpPortStr := config.AppConfig.SMTPPort  //"587"
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Email Verification Code to "+firstName)
	m.SetBody("text/html", fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Email Verification</title>
		</head>
		<body>
			<h2>Hello %s,</h2>
			<p>Your verification code is:</p>
			<h1>%s</h1>
			<p>Please use this code to verify your email address.</p>
		</body>
		</html>
	`, firstName, code))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendPasswordResetEmail(email, token string) error {
  //TODO: Implement the email sending logic to .env variables
	from := config.AppConfig.FromEmail          //"your-email@example.com"
	smtpPass := config.AppConfig.SMTPPass //"your-email-password"
	smtpUser := config.AppConfig.SMTPUser //"your-email-username"
	to := email
	smtpHost := config.AppConfig.SMTPHost    //"smtp.gmail.com"
	smtpPortStr := config.AppConfig.SMTPPort //"587"
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Password Reset")
	m.SetBody("text/html", fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Password Reset</title>
		</head>
		<body>
			<h2>Click the button below to reset your password:</h2>
			<p>
				<a href="http://localhost:3000/password-reset?token=%s" style="display: inline-block; padding: 10px 20px; font-size: 16px; color: #ffffff; background-color: #007bff; text-decoration: none; border-radius: 5px;">Reset Password</a>
			</p>
		</body>
		</html>
	`, token))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// generateEmailVerificationCode generates a random 6-digit code
func GenerateEmailVerificationCode() string {
	return string(randstr.String(20))
}
