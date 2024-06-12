package confirmation

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"time"
)

const (
	testConfirmationCode         = 111111
	confirmationCodeRandMaxLimit = 888889
)

type SMTPClient interface {
	SendEmail(to, subject, body string) error
}

type CodeManager struct {
	smtpClient SMTPClient
	isTest     bool
	r          *rand.Rand
}

func NewCodeManager(smtpClient SMTPClient, isTest bool) *CodeManager {
	return &CodeManager{
		smtpClient: smtpClient,
		isTest:     isTest,
		r:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m *CodeManager) GenerateCode() string {
	if m.isTest {
		return strconv.Itoa(testConfirmationCode)
	}
	return strconv.Itoa(m.r.Intn(confirmationCodeRandMaxLimit) + testConfirmationCode)
}

func (m *CodeManager) SendCode(_ context.Context, opID int64, email, code, langCode string) error {
	if email == "" {
		return errors.New("got empty email")
	}

	message := fmt.Sprintf("Код подтверждения для вывода средств по операции %v: %v", opID, code)
	subject := "Код подтверждения вывода средств Ecomway"

	if langCode == "en" {
		message = fmt.Sprintf("Payout confirmation code for operation %v: %v", opID, code)
		subject = "Ecomway payout confirmation code"
	}

	if !m.isTest {
		if err := m.smtpClient.SendEmail(email, subject, message); err != nil {
			return err
		}
	}

	slog.Info("confirmation code send",
		"email", email,
		"code", code,
		"message", message,
		"operation_id", opID,
	)

	return nil
}
