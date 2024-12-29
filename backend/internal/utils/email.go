package utils

import (
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
)

func VerifyEmail(c *gin.Context) {
    code := c.Param("code")
    // Use the 'code' variable as needed
    // ...existing code...
	fmt.Println(code)
}

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
func SendVerificationCode(email, code string) error {
	from := os.Getenv("FROM") //"your-email@example.com"
	password := os.Getenv("EMAIL_PASSWORD") //"your-email-password"
	to := email
	smtpHost := os.Getenv("SMTP_HOST") //"smtp.gmail.com"
	smtpPort := os.Getenv("SMTP_PORT") //"587"
	
	msg := "From: " + from + "\n" +
	"To: " + to + "\n" +
	"Subject: Email Verification Code\n\n" +
	"Your verification code is: " + code
	
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

// VerifyEmailCode checks if the provided code matches the generated code
func VerifyEmailCode(inputCode, actualCode string) bool {
	return inputCode == actualCode
}

//generateEmailVerificationCode generates a random 6-digit code
func GenerateEmailVerificationCode() string {
	return Encode(randstr.String(20))
}