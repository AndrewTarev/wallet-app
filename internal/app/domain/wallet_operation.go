package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"wallet-app/internal/app/app_errors"
)

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type WalletOperation struct {
	WalletID      uuid.UUID     `json:"walletId" validate:"required"`
	OperationType OperationType `json:"operationType" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        string        `json:"amount" validate:"required,numeric"`
}

var NewValidate = validator.New()

func (op *WalletOperation) Validate() error {
	if err := NewValidate.Struct(op); err != nil {
		return err
	}

	amount, err := op.ParseAmount()
	if err != nil {
		return app_errors.ErrInvalidAmount
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return app_errors.ErrAmountMustBePositive
	}

	return nil
}

func (op *WalletOperation) ParseAmount() (decimal.Decimal, error) {
	return decimal.NewFromString(op.Amount)
}

// GetSignedAmount определяет знак суммы в зависимости от типа операции (DEPOSIT или WITHDRAW)
func (op *WalletOperation) GetSignedAmount() (decimal.Decimal, error) {
	// Преобразуем строку Amount в decimal.Decimal
	amount, err := op.ParseAmount()
	if err != nil {
		return decimal.Zero, err // Если не удалось распарсить — возвращаем ошибку
	}

	// Если операция — снятие средств, делаем сумму отрицательной
	if op.OperationType == Withdraw {
		return amount.Neg(), nil
	}

	// Если операция — пополнение, возвращаем сумму как есть
	return amount, nil
}

func ParseValidationErrors(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errorMessages []string
		for _, ve := range validationErrors {
			field := ve.Field() // Имя поля
			tag := ve.Tag()     // Тег валидации, например, "required" или "oneof"
			switch tag {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", field))
			case "uuid4":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be a valid UUID", field))
			case "oneof":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be one of: %s", field, ve.Param()))
			case "numeric":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be a numeric value", field))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("%s is invalid", field))
			}
		}
		return strings.Join(errorMessages, "; ")
	}
	return "Invalid input data"
}
