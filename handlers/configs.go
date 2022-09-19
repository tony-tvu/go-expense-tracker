package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/cache"
	"github.com/tony-tvu/goexpense/entity"
	"gorm.io/gorm"
)

type ConfigsHandler struct {
	Db    *gorm.DB
	Cache *cache.Configs
}

func (h *ConfigsHandler) RegistrationEnabled(c *gin.Context) {
	configs, err := h.Cache.GetConfigs()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"registration_enabled": configs.RegistrationEnabled,
	})
}

func (h *ConfigsHandler) GetConfigs(c *gin.Context) {
	if _, userType, err := auth.AuthorizeUser(c, h.Db); err != nil || *userType != string(entity.AdminUser) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	configs, err := h.Cache.GetConfigs()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, configs)
}
