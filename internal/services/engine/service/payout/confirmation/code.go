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

type CodeManager struct {
	isTest bool
	r      *rand.Rand
}

func NewCodeManager(isTest bool) *CodeManager {
	return &CodeManager{
		isTest: isTest,
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
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
	if langCode == "en" {
		message = fmt.Sprintf("Withdrawal confirmation code for operation %v: %v", opID, code)
	}

	// TODO: реализовать отправку кода подтверждения на почту
	slog.Info("confirmation code send",
		"email", email,
		"code", code,
		"message", message,
		"operation_id", opID,
	)

	return nil
}
