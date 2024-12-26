package api

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"net"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func VerifyEmail(c *gin.Context) {
    code := c.Param("code")
    // Use the 'code' variable as needed
    // ...existing code...
	fmt.Println(code)
}

// VerifyEmailExistence checks if the email format is valid and if the domain has MX records
func VerifyEmailExistence(email string) (bool, error) {
    // Check email format
    emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(emailRegex)
    if !re.MatchString(email) {
        return false, errors.New("invalid email format")
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
	from := "your-email@example.com"
	password := "your-email-password"
	to := email
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	
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
	return strconv.Itoa(int(math.Floor(100000 + rand.Float64()*(999999-100000))))
}