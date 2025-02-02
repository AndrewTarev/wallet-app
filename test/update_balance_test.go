package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	delivery "wallet-app/internal/app/delivery/http"
	"wallet-app/internal/app/domain"
	"wallet-app/internal/app/services"
	"wallet-app/internal/app/services/mocks"
)

func TestChangeBalance_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок-сервис
	mockService := mocks.NewMockWallet(ctrl)

	// Статичный UUID для теста
	walletID := uuid.Must(uuid.Parse("5979943c-4b3d-4313-9822-b12b47333916"))

	// Ожидаем, что метод ProcessOperation будет вызван с нужной операцией
	mockService.EXPECT().ProcessOperation(gomock.Any(), gomock.Eq(domain.WalletOperation{
		WalletID:      walletID,
		OperationType: domain.OperationType("DEPOSIT"), // Убедитесь, что "DEPOSIT" корректен
		Amount:        "100.00",
	})).Return(nil).Times(1)

	// Создаем сервис с мок-сервисом
	service := &services.Service{
		Wallet: mockService,
	}

	// Создаем хэндлер и роутер
	h := delivery.NewHandler(service)
	router := gin.Default()
	router.POST("/api/v1/wallets/change-balance", h.ChangeBalance)

	// Отправляем запрос с корректными данными
	validOperation := map[string]string{
		"walletId":      walletID.String(), // Используем статичный UUID
		"operationType": "DEPOSIT",         // Используем правильный тип операции
		"amount":        "100.00",
	}

	requestBody, err := json.Marshal(validOperation)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/api/v1/wallets/change-balance", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Проверяем, что вернулся статус OK
	assert.Equal(t, http.StatusOK, resp.Code)

	// Проверяем тело ответа
	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Operation completed", response["message"])
}

func TestChangeBalance_InvalidRequestFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок-сервис
	mockService := mocks.NewMockWallet(ctrl)

	// Создаем сервис с мок-сервисом
	service := &services.Service{
		Wallet: mockService,
	}

	// Создаем хэндлер и роутер
	h := delivery.NewHandler(service)
	router := gin.Default()
	router.POST("/api/v1/wallets/change-balance", h.ChangeBalance)

	// Отправляем некорректный запрос (например, без тела)
	req, _ := http.NewRequest("POST", "/api/v1/wallets/change-balance", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Проверяем, что вернулась ошибка с статусом BadRequest
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestChangeBalance_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок-сервис
	mockService := mocks.NewMockWallet(ctrl)

	// Создаем сервис с мок-сервисом
	service := &services.Service{
		Wallet: mockService,
	}

	// Создаем хэндлер и роутер
	h := delivery.NewHandler(service)
	router := gin.Default()
	router.POST("/api/v1/wallets/change-balance", h.ChangeBalance)

	// Отправляем запрос с некорректным walletId (невалидный UUID)
	invalidOperation := map[string]string{
		"walletId":      "invalid-uuid", // Некорректный UUID
		"operationType": "debit",
		"amount":        "100.00",
	}

	requestBody, err := json.Marshal(invalidOperation)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/api/v1/wallets/change-balance", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Проверяем, что вернулась ошибка с статусом BadRequest
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestChangeBalance_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок-сервис
	mockService := mocks.NewMockWallet(ctrl)

	// Создаем сервис с мок-сервисом
	service := &services.Service{
		Wallet: mockService,
	}

	// Создаем хэндлер и роутер
	h := delivery.NewHandler(service)
	router := gin.Default()
	router.POST("/api/v1/wallets/change-balance", h.ChangeBalance)

	// Отправляем запрос с некорректными данными операции (например, пустым amount)
	invalidOperation := domain.WalletOperation{
		WalletID:      uuid.New(),
		OperationType: "debit",
		Amount:        "",
	}

	requestBody, err := json.Marshal(invalidOperation)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/api/v1/wallets/change-balance", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Проверяем, что вернулась ошибка с статусом BadRequest
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
