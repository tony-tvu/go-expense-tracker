package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/cache"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
)

type ConfigsHandler struct {
	Db           *database.MongoDb
	ConfigsCache *cache.Configs
}

func (h *ConfigsHandler) RegistrationEnabled(c *gin.Context) {
	configs, err := h.ConfigsCache.GetConfigs()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"registration_enabled": configs.RegistrationEnabled,
	})
}

func (h *ConfigsHandler) GetConfigs(c *gin.Context) {
	if _, userType, err := auth.AuthorizeUser(c, h.Db); err != nil || *userType != models.AdminUser {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	configs, err := h.ConfigsCache.GetConfigs()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, configs)
}
