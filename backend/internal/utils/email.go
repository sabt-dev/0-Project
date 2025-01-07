package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

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
	from := os.Getenv("FROM")          //"your-email@example.com"
	smtpPass := os.Getenv("SMTP_PASS") //"your-email-password"
	smtpUser := os.Getenv("SMTP_USER") //"your-email-username"
	to := email
	smtpHost := os.Getenv("SMTP_HOST")    //"smtp.gmail.com"
	smtpPortStr := os.Getenv("SMTP_PORT") //"587"
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Email Verification Code to "+firstName)
	m.SetBody("text/plain", "Your verification code is: "+code)

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
	from := os.Getenv("FROM")          //"your-email@example.com"
	smtpPass := os.Getenv("SMTP_PASS") //"your-email-password"
	smtpUser := os.Getenv("SMTP_USER") //"your-email-username"
	to := email
	smtpHost := os.Getenv("SMTP_HOST")    //"smtp.gmail.com"
	smtpPortStr := os.Getenv("SMTP_PORT") //"587"
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Password Reset")
	m.SetBody("text/plain",fmt.Sprintf("Click the link to reset your password: http://localhost:3000/password-reset?token=%s", token))

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
