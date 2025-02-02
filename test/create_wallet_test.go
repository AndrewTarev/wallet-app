package test

import (
	"encoding/json"
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

func TestCreateWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок-сервис
	mockService := mocks.NewMockWallet(ctrl)

	// Настроим ожидания на мок-сервис
	mockService.EXPECT().CreateWallet(gomock.Any()).Return(domain.Wallet{
		ID:      uuid.New(),
		Balance: decimal.NewFromInt(0),
	}, nil).Times(1)

	// Создаем сервис, передавая мок-сервис, реализующий интерфейс Wallet
	service := &services.Service{
		Wallet: mockService, // Передаем мок-сервис как реализацию интерфейса Wallet
	}

	// Создаем хэндлер
	h := delivery.NewHandler(service) // Передаем сервис в хэндлер
	router := gin.Default()
	router.POST("/api/v1/create-wallet", h.CreateWallet)

	// Отправляем запрос
	req, _ := http.NewRequest("POST", "/api/v1/create-wallet", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Проверяем ответ
	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "walletId")
	assert.Contains(t, response, "balance")
}
