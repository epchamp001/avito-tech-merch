package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Service interface {
	GetInfo(c *gin.Context)
	SendCoin(c *gin.Context)
	BuyItem(c *gin.Context)
	Auth(c *gin.Context)
}

// это пока временно и тут надо рефакторить, будет другая НЕПУСТАЯ структура
type userService struct{}

// NewUserService создает новый экземпляр сервиса
func NewUserService() Service {
	return &userService{}
}

// GetInfo обработчик /api/info
func (s *userService) GetInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Информация о пользователе",
	})
}

// SendCoin обработчик /api/sendCoin
func (s *userService) SendCoin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Монеты успешно отправлены",
	})
}

// BuyItem обработчик /api/buy/:item
func (s *userService) BuyItem(c *gin.Context) {
	item := c.Param("item") // Получаем item из URL
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Покупка завершена",
		"item":    item,
	})
}

// Auth обработчик /api/auth
func (s *userService) Auth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Авторизация прошла успешно",
	})
}
