package test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	delivery "wallet-app/internal/app/delivery/http"
	"wallet-app/internal/app/domain"
	"wallet-app/internal/app/services"
	"wallet-app/internal/app/services/mocks"
)

func TestGetBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Генерируем новый ID кошелька
	walletID := uuid.New()

	// Создаем мок-сервис
	mockWallet := mocks.NewMockWallet(ctrl)

	// Ожидаем, что метод GetBalance будет вызван с walletID и вернет нужный баланс
	expectedBalance := decimal.NewFromInt(100)
	mockWallet.
		EXPECT().
		GetBalance(gomock.Any(), gomock.Eq(walletID)). // Ожидаем вызов именно с этим UUID
		Return(expectedBalance, nil).Times(1)

	// Оборачиваем мок-сервис в сервис
	service := &services.Service{Wallet: mockWallet}

	// Создаем хэндлер
	h := delivery.NewHandler(service)
	router := gin.Default()
	router.GET("/api/v1/wallets/:walletId", h.GetBalance)

	// Отправляем запрос с нужным ID кошелька
	req, _ := http.NewRequest("GET", "/api/v1/wallets/"+walletID.String(), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Проверяем статус-код ответа
	assert.Equal(t, http.StatusOK, resp.Code)

	// Проверяем содержимое ответа
	var response domain.WalletBalance
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, response.Balance)
}

func TestGetBalance_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Генерация UUID для теста
	walletID := uuid.New()

	// Создаем мок-сервис
	mockService := mocks.NewMockWallet(ctrl)

	// Ожидаем, что метод GetBalance будет вызван с UUID и вернет ошибку
	mockService.EXPECT().GetBalance(gomock.Any(), walletID).Return(decimal.Decimal{}, errors.New("wallet not found")).Times(1)

	// Создаем сервис с мок-сервисом
	service := &services.Service{
		Wallet: mockService,
	}

	// Создаем хэндлер
	h := delivery.NewHandler(service)
	router := gin.Default()
	router.GET("/api/v1/wallets/:walletId", h.GetBalance)

	// Отправляем запрос с правильным ID кошелька
	req, _ := http.NewRequest("GET", "/api/v1/wallets/"+walletID.String(), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Проверяем, что возвращен статус NotFound
	assert.Equal(t, http.StatusNotFound, resp.Code)

	// Проверяем, что в ответе указана ошибка
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Wallet not found", response["error"])
}

func TestGetBalance_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем хэндлер
	h := delivery.NewHandler(nil)
	router := gin.Default()
	router.GET("/api/v1/wallets/:walletId", h.GetBalance)

	// Отправляем запрос с некорректным UUID
	req, _ := http.NewRequest("GET", "/api/v1/wallets/invalid-uuid", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Проверяем, что возвращена ошибка с корректным статусом
	assert.Equal(t, http.StatusForbidden, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid UUID format", response["error"])
}
