package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wallet-app/internal/app/domain"
)

// ChangeBalance обрабатывает изменение баланса кошелька.
//
// @Summary Изменение баланса кошелька
// @Description Пополнение или снятие средств с кошелька
// @Tags wallets
// @Accept json
// @Produce json
// @Param request body domain.WalletOperation true "Данные операции"
// @Success 200 {object} SuccessResponse "Операция выполнена"
// @Failure 400 {object} ErrorResponse "Ошибка валидации данных"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /wallet [post]
func (h *Handler) ChangeBalance(c *gin.Context) {
	var op domain.WalletOperation

	// Привязываем тело запроса к структуре
	if err := c.ShouldBindJSON(&op); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), "Invalid request format")
		return
	}

	// Валидация данных операции
	if err := op.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	// Обрабатываем операцию (пополнение или снятие)
	if err := h.services.ProcessOperation(c.Request.Context(), op); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), "Error in processing the operation")
		return
	}

	response := SuccessResponse{Message: "Operation completed"}

	c.JSON(http.StatusOK, response)
}

// GetBalance получает текущий баланс кошелька.
//
// @Summary Получение баланса кошелька
// @Description Возвращает баланс указанного кошелька
// @Tags wallets
// @Accept json
// @Produce json
// @Param walletId path string true "UUID кошелька"
// @Success 200 {object} domain.WalletBalance "Баланс кошелька"
// @Failure 400 {object} ErrorResponse "Неверный UUID"
// @Failure 404 {object} ErrorResponse "Кошелек не найден"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /wallets/{walletId} [get]
func (h *Handler) GetBalance(c *gin.Context) {
	walletID := c.Param("walletId")
	walletUUID, err := uuid.Parse(walletID)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error(), "Invalid UUID format")
		return
	}

	balance, err := h.services.GetBalance(c.Request.Context(), walletUUID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error(), "Wallet not found")
		return
	}

	response := domain.WalletBalance{
		Balance: balance,
	}

	c.JSON(http.StatusOK, response)
}

// CreateWallet создает новый кошелек.
//
// @Summary Создание нового кошелька
// @Description Генерирует новый кошелек с начальным балансом 0
// @Tags wallets
// @Accept json
// @Produce json
// @Success 201 {object} domain.Wallet "Созданный кошелек"
// @Failure 400 {object} ErrorResponse "Ошибка при создании"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /create-wallet [post]
func (h *Handler) CreateWallet(c *gin.Context) {
	wallet, err := h.services.CreateWallet(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), "Wallet creation error")
		return
	}

	c.JSON(http.StatusCreated, wallet)
}
